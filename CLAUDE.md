# Rewinds — CLAUDE.md

## App concept
A music taste platform where users rate albums and artists via an intelligent comparison flow (Beli-style), vote head-to-head in ELO matchups, discover underground artists, make leaderboard rank predictions using a points-based betting system, get tour alerts for artists they've rated, and share their taste via customizable public profiles and shareable graphics.

**Domain:** rewinds.me

---

## Tech stack
| Layer | Technology |
|---|---|
| Mobile | React Native + Expo + Expo Router |
| Styling | NativeWind |
| Backend | Go + Chi router |
| Database | Supabase (Postgres hosting only — all queries via jackc/pgx direct connection) |
| Cache / Leaderboards | Redis (Railway plugin) |
| File storage | Cloudflare R2 |
| Auth | Spotify OAuth (primary) + email/password fallback — handled entirely by Go backend, JWT issued for all sessions |
| Music data | Spotify API (catalog, OAuth, top artists, recently played, ongoing sync) |
| Tour / ticket data | Songkick API + Ticketmaster Discovery API |
| Payments | Stripe (web checkout, US) + StoreKit (non-US) via RevenueCat |
| Backend deployment | Railway |
| Profile pages (web) | Next.js + Vercel |
| CI/CD | GitHub Actions |
| Mobile builds | Expo EAS Build |

## Architecture
```
React Native → Go API (REST) → Postgres (direct via jackc/pgx)
```
- React Native never connects to Supabase or Postgres directly
- React Native never uses a Supabase client library
- All data access is gated through the Go API
- Supabase is infrastructure only — it hosts Postgres and provides the dashboard
- Go backend connects to Postgres via a single SUPABASE_DB_URL connection string
- No Supabase client library used anywhere in the codebase

---

## Database schema

### users
```
id                    uuid PK
username              text unique
email                 text unique (nullable for Spotify-only signups)
avatar_url            text
bio                   text
is_premium            bool default false
credits               int default 10
credits_reset_at      timestamp
points                int default 1000        -- current balance, not lifetime. starts at 1000 signup bonus
points_rank           int                     -- global rank, recomputed hourly from Redis
spotify_id            text unique (nullable)  -- null for email signups
spotify_access_token  text                    -- encrypted, refreshed periodically
spotify_refresh_token text                    -- encrypted
spotify_token_expiry  timestamp
onboarding_complete   bool default false      -- true once user has rated ≥5 items (for comparison flow) or skipped
created_at            timestamp
```

### artists
```
id            uuid PK
spotify_id    text unique
name          text
image_url     text
genres        text[]
followers     int
popularity    int
created_at    timestamp
```

### elo_scores
```
id            uuid PK
artist_id     uuid FK -> artists
genre         text
score         int default 1000
vote_count    int default 0
updated_at    timestamp
```

### votes
```
id            uuid PK
user_id       uuid FK -> users
winner_id     uuid FK -> artists
loser_id      uuid FK -> artists
genre         text
created_at    timestamp
```

### ratings
```
id                uuid PK
user_id           uuid FK -> users
artist_id         uuid FK -> artists
album_id          uuid FK -> albums (nullable)
track_id          uuid FK -> tracks (nullable)
score             decimal(3,1)
rating_method     text (comparison | manual)
annotation        text    -- personal note shown on profile next to the rating ("why I gave this 8.2")
created_at        timestamp
```

### rating_comparisons
```
id            uuid PK
user_id       uuid FK -> users
item_a_id     uuid
item_b_id     uuid
item_a_type   text (artist | album | track)
winner_id     uuid
created_at    timestamp
```

### comments
```
id            uuid PK
user_id       uuid FK -> users
body          text
page_type     text (artist | album | track)   -- for page-level comments
page_id       uuid                             -- FK to artists, albums, or tracks
rating_id     uuid FK -> ratings (nullable)    -- set only for rating-level comments
parent_id     uuid FK -> comments (nullable)   -- for threaded replies on page comments
created_at    timestamp
updated_at    timestamp
```
> Two comment types share this table:
> - **Page-level**: `page_type` + `page_id` set, `rating_id` null. Public discussion on artist/album/track pages. Supports replies via `parent_id`.
> - **Rating annotation**: `rating_id` set, visible on the user's profile next to their rating. Personal note ("why I gave this 8.2"). No replies on annotations.

### folders
```
id            uuid PK
user_id       uuid FK -> users
name          text
description   text (nullable)
is_public     bool default false
created_at    timestamp
updated_at    timestamp
```

### folder_ratings
```
id            uuid PK
folder_id     uuid FK -> folders
rating_id     uuid FK -> ratings
added_at      timestamp
```
> Junction table — a rating can live in multiple folders simultaneously.

### albums
```
id            uuid PK
spotify_id    text unique
artist_id     uuid FK -> artists
name          text
image_url     text
release_date  date
created_at    timestamp
```

### tracks
```
id            uuid PK
spotify_id    text unique
album_id      uuid FK -> albums
artist_id     uuid FK -> artists
name          text
duration_ms   int
created_at    timestamp
```

### predictions
```
id                uuid PK
user_id           uuid FK -> users
artist_id         uuid FK -> artists
genre             text
current_rank      int            -- artist's rank at time of bet
target_rank       int            -- rank user predicts artist will reach or surpass
deadline          timestamp
wager             int            -- points deducted immediately at bet time
payout_multiplier decimal(4,2)  -- computed at bet time based on climb size
outcome           text (pending | won | lost)
predicted_at      timestamp
resolved_at       timestamp
```

### point_transactions
```
id            uuid PK
user_id       uuid FK -> users
amount        int        -- positive = earned, negative = spent/lost
reason        text       -- quest_daily | quest_weekly | promo_rating | prediction_won | prediction_lost | prediction_wager | signup_bonus
reference_id  uuid       -- nullable FK to votes, ratings, or predictions
created_at    timestamp
```

### spotify_sync_log
```
id            uuid PK
user_id       uuid FK -> users
synced_at     timestamp
top_artists   jsonb    -- raw snapshot of top artists at sync time
recent_tracks jsonb    -- raw snapshot of recently played at sync time
```
> Keeps a history of Spotify syncs. Used to suggest new items to rate when taste changes over time.

### tour_alerts
```
id            uuid PK
user_id       uuid FK -> users
artist_id     uuid FK -> artists
enabled       bool default true
created_at    timestamp
```

### tour_events
```
id              uuid PK
artist_id       uuid FK -> artists
source          text (songkick | ticketmaster)
external_id     text unique
venue_name      text
city            text
country         text
event_date      timestamp
ticket_url      text
created_at      timestamp
```

### promo_campaigns
```
id                uuid PK
artist_id         uuid FK -> artists
label_name        text
budget_pts        int         -- total points pool to distribute to users
pts_per_rating    int         -- points awarded per completed promo rating
ratings_target    int         -- campaign ends when this many ratings are collected
ratings_count     int default 0
genre             text        -- target genre for user matching
status            text (active | paused | completed)
created_at        timestamp
ends_at           timestamp
```

### promo_ratings
```
id              uuid PK
campaign_id     uuid FK -> promo_campaigns
user_id         uuid FK -> users
rating_id       uuid FK -> ratings
previewed       bool default false   -- did user play the preview before rating
pts_awarded     int
created_at      timestamp
```
> One row per user per campaign. Prevents claiming points twice on same campaign.

### quests
```
id              uuid PK
type            text (daily | weekly)
title           text
description     text
pts_reward_free     int
pts_reward_premium  int     -- always higher than free
action          text        -- vote | rate | comment | promo_rating | prediction | share | discover
target_count    int         -- how many times action must be completed
is_active       bool default true
created_at      timestamp
```

### user_quest_progress
```
id              uuid PK
user_id         uuid FK -> users
quest_id        uuid FK -> quests
progress        int default 0       -- current count toward target_count
completed       bool default false
pts_claimed     bool default false
period_start    timestamp           -- which daily/weekly period this belongs to
created_at      timestamp
```
> Reset daily quests at midnight, weekly quests on Monday midnight.
> New rows created each period — don't update old ones.

### profile_customizations
```
id              uuid PK
user_id         uuid FK -> users
theme_color     text
layout          text
pinned_artists  uuid[]
mount_rushmore  uuid[4]
created_at      timestamp
```

---

## Folder structure

### Go backend
```
/backend
  /cmd
    main.go
  /internal
    /handlers
      auth.go
      spotify_auth.go
      artists.go
      albums.go
      tracks.go
      votes.go
      ratings.go
      comparisons.go
      leaderboard.go
      profiles.go
      predictions.go
      points.go
      comments.go
      folders.go
      quests.go
      promo.go
      tours.go
    /services
      elo.go
      matchmaking.go
      credits.go
      spotify.go
      spotify_sync.go
      graphics.go
      comparison_rating.go
      points.go
      rankings.go
      quests.go
      promo.go
      tour_alerts.go
      predictions.go
    /models
      user.go
      artist.go
      album.go
      track.go
      vote.go
      rating.go
      comparison.go
      prediction.go
      point_transaction.go
      comment.go
      folder.go
      quest.go
      promo_campaign.go
      tour_event.go
    /middleware
      auth.go
      ratelimit.go
    /db
      supabase.go
      redis.go
  /config
    config.go
  go.mod
  go.sum
  Dockerfile
```

### React Native
```
/mobile
  /app
    /(tabs)
      index.tsx        (voting screen)
      leaderboard.tsx  (artist leaderboard + user rank tab)
      discover.tsx
      profile.tsx
    /quests/
      index.tsx        (daily + weekly quest board)
    /promo/
      [id].tsx         (promo rating flow for a campaign)
    /folders/
      index.tsx        (all folders)
      [id].tsx         (folder contents)
      new.tsx          (create folder)
    /spotify/
      playlists.tsx    (list user's Spotify playlists)
      [id].tsx         (playlist contents with rate buttons)
    /artist/[id].tsx
    /album/[id].tsx
    /predictions/
      index.tsx        (open + resolved predictions)
      new.tsx          (place a new prediction)
    /rate/
      index.tsx        (choose rating method)
      comparison.tsx   (Beli-style flow)
      manual.tsx       (manual score entry)
    /tours/[artistId].tsx
    /auth/
      login.tsx
      signup.tsx
      spotify-callback.tsx
    /onboarding/
      index.tsx        (rate these first queue)
    /premium/
      index.tsx
  /components
    /ui
      Button.tsx
      Card.tsx
      ArtistCard.tsx
    VotingCard.tsx
    LeaderboardRow.tsx
    RatingStars.tsx
    GraphicExport.tsx
    ComparisonCard.tsx
    TourAlertBanner.tsx
    PredictionCard.tsx
    PointsDisplay.tsx
    QuestCard.tsx
    PromoCard.tsx
    CommentThread.tsx
    CommentInput.tsx
    FolderCard.tsx
    PlaylistTrackRow.tsx
  /hooks
    useVoting.ts
    useCredits.ts
    useProfile.ts
    useRating.ts
    useTourAlerts.ts
    usePredictions.ts
    usePoints.ts
    useAuth.ts
    useOnboarding.ts
    useComments.ts
    useFolders.ts
    useQuests.ts
    usePromo.ts
    useSpotifyPlaylists.ts
  /lib
    api.ts          -- all Go API calls, sets Authorization header with JWT
    storage.ts      -- secure local JWT storage (expo-secure-store)
    revenuecat.ts   -- RevenueCat SDK (payments only, no DB access)
  /constants
    colors.ts
    genres.ts
  app.json
  package.json
```

---

## Auth flow

**Primary: Spotify OAuth**
- User taps "Continue with Spotify"
- Backend exchanges Spotify auth code for access + refresh tokens
- Tokens stored encrypted on user record in Postgres
- Pull `user-top-read` + `user-read-recently-played` scopes immediately on signup
- Go backend creates user record, issues internal JWT for all subsequent app sessions

**Fallback: Email/password**
- Go backend handles registration and login directly — no Supabase Auth
- Passwords hashed with bcrypt in Go before storing in Postgres
- Go backend issues JWT on successful login
- No Spotify data on signup — user starts with empty taste profile
- Can connect Spotify later from profile settings to unlock onboarding queue retroactively

**JWT:**
- All sessions (Spotify and email) use a JWT issued by the Go backend
- React Native stores JWT in secure local storage and sends it as `Authorization: Bearer <token>` on every request
- Go auth middleware validates JWT on every protected route

**Ongoing Spotify sync (weekly cron):**
- Re-pull top artists + recently played for all users with connected Spotify
- Compare against existing ratings — surface unrated artists as "new suggestions"
- Store snapshot in spotify_sync_log

**Spotify playlist access:**
- Scopes: `user-top-read`, `user-read-recently-played`, `playlist-read-private`, `playlist-read-collaborative`
- User can browse their own Spotify playlists inside Rewinds
- Each track shows rated/unrated status — green badge if rated, empty if not
- "Rate all unrated" queues unrated tracks for quick sequential rating
- Playlists are read-only — Rewinds never writes to or modifies Spotify playlists

---

## Onboarding flow

After auth (Spotify path):
1. Pull user's top 50 artists + recently played from Spotify
2. Filter to artists not yet rated by user
3. Show "Rate these first" queue — "We found X artists you know. Rate a few to calibrate your taste."
4. User rates items one by one:
   - **Manual rating** (default during onboarding — just a star tap, fast)
   - No forced comparison flow yet — not enough anchors
5. After rating 5 items → `onboarding_complete = true`, comparison flow unlocks
6. User enters main app

After auth (email path):
1. Skip directly to main app — no onboarding queue
2. Profile shows "Connect Spotify to get personalized suggestions" prompt
3. Comparison flow shows manual rating only until 5 ratings exist

**Comparison flow gate:**
- `rating_count < 5` → always use manual rating, hide comparison option
- `rating_count >= 5` → comparison flow available, shown as default option

---



```go
func calculateELO(winnerScore, loserScore, kFactor int) (int, int) {
    expectedWinner := 1.0 / (1.0 + math.Pow(10, float64(loserScore-winnerScore)/400))
    expectedLoser := 1.0 - expectedWinner

    newWinner := winnerScore + int(float64(kFactor) * (1 - expectedWinner))
    newLoser := loserScore + int(float64(kFactor) * (0 - expectedLoser))

    return newWinner, newLoser
}
// K-factor: 32 for new artists (<100 votes), 16 for established (100+ votes)
```

---

## Matchmaking algorithm

```go
func getMatchup(genre string, upsetProbability float64) (Artist, Artist) {
    if rand.Float64() < upsetProbability {
        // Upset: one from top 20%, one from bottom 20%
        top := getTopArtists(genre, 20)
        bottom := getBottomArtists(genre, 20)
        return randomFrom(top), randomFrom(bottom)
    }
    // Normal: two artists within 100 ELO points of each other
    anchor := getRandomArtist(genre)
    opponent := getArtistNearELO(genre, anchor.ELO, 100)
    return anchor, opponent
}
```

---

## Rating flow

Users always choose their rating method — the comparison flow is **never forced**.

### Method A — Manual (always available)
User enters a decimal score directly (e.g. `7.4`) or uses a star/slider input.
Saved immediately as `rating_method: "manual"`.

### Method B — Comparison flow (available once ≥ 5 ratings exist)

**Step 1 — Category selection**

User picks one of 5 broad sentiment categories:

| Label | Score range | Starting score (first in category) |
|---|---|---|
| Very Bad | 0.0 – 2.0 | 1.0 |
| Bad | 2.0 – 4.0 | 3.0 |
| Ok | 4.0 – 6.0 | 5.0 |
| Good | 6.0 – 8.0 | 7.0 |
| Amazing | 8.0 – 10.0 | 9.0 |

If this is the first item in the chosen category → skip to Step 3, assign starting score directly.

**Step 2 — Server-driven binary search (within category)**

- Session state (`low`, `high`, snapshot of that category's rated items) stored in Redis
- `GET /comparisons/next` returns the midpoint item to compare against
- User answers "better" or "worse" → `POST /comparisons` records the result and advances the session
- Repeat until `low > high` — position falls out as the final `low` value
- Boundary check: if the song lands at the top edge of its category, one extra comparison against the lowest-ranked item in the category above confirms or recategorizes it (same at bottom edge)

**Step 3 — Score calculation**

Manual ratings act as fixed window boundaries. Comparison items are distributed evenly within the window defined by the nearest manual scores above and below (or category bounds if no manual exists).

```go
// ScoreInWindow distributes comparison items evenly within [lower, upper].
// Scores approach but never reach `upper` — so comparison scores never equal
// a manual score ranked above them, and 10.0 is only reachable via manual entry.
//
// position: 0-indexed (0 = best in window)
// nTotal:   total comparison items in this window including this one
// lower:    nearest manual score below (or category minimum)
// upper:    nearest manual score above (or category maximum)
func ScoreInWindow(position, nTotal int, lower, upper float64) float64 {
    if nTotal <= 1 {
        return (upper + lower) / 2.0 // midpoint when only one item in window
    }
    step := (upper - lower) / float64(nTotal+1)
    return upper - step*float64(position+1)
}
```

**Example** — 1 manual at 9.5, 1 comparison above it:
```
window [9.5, 10.0], nTotal=1 → (9.5 + 10.0) / 2 = 9.75

Result: 9.75 (comparison) | 9.5 (manual)
```

**Example** — 1 manual at 9.5, 1 comparison above, 1 below:
```
upper window [9.5, 10.0], nTotal=1 → 9.75 (comparison)
lower window [8.0, 9.5],  nTotal=1 → 8.75 (comparison)

Result: 9.75 (comparison) | 9.5 (manual) | 8.75 (comparison)
```

**Step 4 — Confirm or override**

User sees the derived score and chooses:
- **Accept** → saved as `rating_method: "comparison"`
- **Edit the score** → saved as `rating_method: "manual"` (user's decimal wins; score is never auto-recalculated)
- **Discard** → nothing saved

**Step 5 — Recompute category scores after insertion**

After a new `"comparison"` rating is saved, all `"comparison"` ratings in that category are rescored. The full sorted list (manual + comparison together) is used to identify windows, but only comparison scores are updated.

```
Sort all items in category by score descending.
Split into windows at each manual rating:

  [categoryMax] ... manual(9.5) ... manual(8.0) ... [categoryMin]
       window 1         window 2

For each window:
  compItems = comparison items in this window
  For each at position i:
      score = ScoreInWindow(i, len(compItems), lower, upper)
```

Manual scores are fixed anchors — they bound the windows but are never recalculated. This prevents any comparison score from ever exceeding a manual score ranked above it.

---

## Points system

**Points is current balance** — goes up when earned, down when lost on bets. Rank is derived from current points so betting is always a real risk/reward decision. Everyone starts with 1,000 points on signup.

**Earning points — quests only (no passive per-action rewards):**

Daily quests (reset midnight):
| Quest | Free | Premium |
|---|---|---|
| Vote in 5 head-to-head matchups | +50 pts | +75 pts |
| Rate 1 new artist or album | +75 pts | +115 pts |
| Listen to a preview and rate | +60 pts | +90 pts |
| Comment on an artist or album page | +40 pts | +60 pts |
| Complete a promo rating | +150 pts | +225 pts |

Weekly quests (reset Monday midnight):
| Quest | Free | Premium |
|---|---|---|
| Rate 5 albums in the same genre | +300 pts | +450 pts |
| Make a prediction that resolves this week | +200 pts | +300 pts |
| Discover and rate 3 artists under 10k followers | +400 pts | +600 pts |
| Rate artists from 3 different genres | +250 pts | +375 pts |
| Share a taste graphic to social | +200 pts | +300 pts |

**Signup bonus:** +1,000 pts (once, on account creation)

**Promo ratings:**
- Always available when active campaigns exist in user's matched genres
- Highest single point reward of any quest — incentivizes engagement with promoted content
- User rates the promoted artist/track (preview optional but shown if available)
- Points awarded on rating submission, logged to point_transactions with reason `promo_rating`
- One claim per user per campaign (enforced via promo_ratings table)

**Spending points:**
- Prediction wagers (deducted immediately at bet time)

**Betting limits:**
- Free: max 1,000 pts per bet, max 10 open predictions
- Premium: max 10,000 pts per bet, max 25 open predictions

**Global rank:** Every user has a rank based on current points balance. Displayed on profile as `Rank #3 · 142,830 pts`. Updates hourly.

---

## Prediction system

Users bet points that an artist will climb to a target leaderboard rank within a genre before a deadline. All ranks are genre-specific.

**Payout multiplier table** (computed at bet time based on rank climb distance):
| Predicted rank climb | Multiplier |
|---|---|
| 1–10 ranks | 1.5x |
| 11–25 ranks | 3x |
| 26–50 ranks | 6x |
| 51–100 ranks | 12x |
| Into top 10 | 20x |

**Resolution (hourly cron job):**
- Pull artist's current rank from Redis
- If current rank ≤ target rank → won, credit `wager × multiplier` to points, log to point_transactions
- If deadline passed and rank not reached → lost (wager already deducted at bet time)
- Push notification either way

**Anti-abuse:**
- Minimum wager: 100 pts
- Artist must have ≥ 50 votes to be eligible for predictions
- Can't have two open predictions on same artist in same genre simultaneously

---

## Tour alerts

When a user rates an artist, automatically create a tour_alert record for that artist.

**Cron job (every 6 hours):**
```
For each artist with tour_alerts enabled:
  Query Songkick API + Ticketmaster API for new events
  Deduplicate against tour_events table by external_id
  For each new event:
    Store in tour_events
    Find users with tour_alert for this artist
    Filter by users within 100mi of event city
    Queue push notification: "🎟️ [Artist] just announced a show near you!"
```

**Ticket links:** Generate SeatGeek affiliate links for ticket purchases. Ticketmaster as fallback. SeatGeek pays ~$5 flat per ticket sold.

**APIs:**
- Songkick API (free, strong indie/underground coverage)
- Ticketmaster Discovery API (free tier, mainstream catalog)
- Use both, deduplicate by `external_id`

---

## Redis data structures

```
# ELO leaderboard per genre (sorted set)
ZADD leaderboard:hiphop <elo_score> <artist_id>
ZREVRANGE leaderboard:hiphop 0 99        # top 100 artists

# Weekly snapshot for "rising" movers calculation
ZADD leaderboard:hiphop:week:<week_num> <elo_score> <artist_id>

# User credits (resets daily at midnight)
SET credits:<user_id> 10
EXPIREAT credits:<user_id> <midnight_unix_timestamp>

# Matchup deduplication (avoid repeating same pair same day)
SADD seen_matchups:<user_id> <artist1_id>:<artist2_id>
EXPIRE seen_matchups:<user_id> 86400

# Global user rankings by current points balance (sorted set)
ZADD user_rankings <points_balance> <user_id>
ZREVRANK user_rankings <user_id>         # 0-indexed, add 1 for display

# Artist rank per genre (same data as leaderboard, named for clarity in prediction resolution)
ZREVRANK leaderboard:hiphop <artist_id>  # current rank for prediction checks
```

---

## Sprint plan

### Sprint 1 — Foundation (Week 1-2)
**Backend**
- [ ] Initialize Go project with Chi router
- [ ] Connect Supabase Postgres client
- [ ] Connect Redis client
- [ ] Implement JWT auth middleware
- [ ] POST /auth/spotify (exchange Spotify code for tokens, create/login user, issue JWT)
- [ ] POST /auth/register (email fallback)
- [ ] POST /auth/login (email fallback)
- [ ] GET /auth/spotify/callback
- [ ] POST /auth/spotify/refresh (refresh Spotify access token)
- [ ] GET /spotify/onboarding (pull top artists + recently played for onboarding queue)
- [ ] GET /spotify/playlists (list user's Spotify playlists)
- [ ] GET /spotify/playlists/:playlist_id/tracks (tracks in a playlist with Rewinds rating status)
- [ ] GET /health
- [ ] Weekly cron: Spotify sync for all connected users → surface new unrated artists

**Mobile**
- [ ] Initialize Expo project with TypeScript
- [ ] Set up Expo Router tab navigation
- [ ] Set up NativeWind
- [ ] Login screen: "Continue with Spotify" (primary CTA) + "Continue with email" (secondary)
- [ ] Spotify OAuth flow (Expo AuthSession)
- [ ] Email signup/login screens (fallback)
- [ ] Onboarding screen: "Rate these first" queue (manual star rating only, 5 item minimum)
- [ ] Skip onboarding option (email users, or Spotify users who don't want to rate yet)
- [ ] Gate comparison flow behind rating_count ≥ 5 check
- [ ] Persist JWT token locally

**Infrastructure**
- [ ] Create GitHub repo (rewinds)
- [ ] Set up Railway project (Go backend + Redis)
- [ ] Set up Supabase project and run schema migrations
- [ ] Set up GitHub Actions for Go tests on push
- [ ] Register Spotify developer app (scopes: user-top-read, user-read-recently-played, playlist-read-private, playlist-read-collaborative)
- [ ] Register Songkick developer account
- [ ] Register SeatGeek affiliate account

---

### Sprint 2 — Core voting loop + points foundation (Week 3-4)
**Backend**
- [ ] GET /artists/matchup?genre=x
- [ ] POST /votes (record vote, update ELO, award +10/+20 pts)
- [ ] Implement ELO calculation service
- [ ] Write ELO scores to Redis sorted sets per genre
- [ ] GET /leaderboard?genre=x
- [ ] GET /leaderboard/rising?genre=x
- [ ] Daily cron: reset credits, award daily login points
- [ ] Spotify API integration (fetch + sync artist catalog)
- [ ] Points service (award, deduct, log to point_transactions)
- [ ] User ranking sorted set in Redis (ZADD user_rankings)
- [ ] GET /leaderboard/users (global rank leaderboard top 100)
- [ ] GET /users/:id/rank (rank + 5 neighbors above and below)

**Mobile**
- [ ] Voting screen (two artist cards, tap to vote)
- [ ] Vote result animation (winner highlight, loser fade)
- [ ] Credits remaining display
- [ ] Leaderboard screen (genre filter + user rank tab)
- [ ] Rising artists tab
- [ ] Spotify 30s preview playback inline while voting
- [ ] Points balance + Rank #X display on profile tab

---

### Sprint 3 — Ratings + comparison flow + comments + folders (Week 5-6)
**Backend**
- [ ] POST /ratings (supports both rating methods, awards +25/+50 pts)
- [ ] GET /ratings?user_id=x
- [ ] PUT /ratings/:id (update score or annotation)
- [ ] DELETE /ratings/:id
- [ ] POST /comparisons (log single binary comparison)
- [ ] GET /comparisons/next?user_id=x&item_id=y
- [ ] Implement binary search comparison rating service
- [ ] GET /artists/:id (with ELO, ratings, vote count, comment count)
- [ ] GET /albums/:id
- [ ] GET /tracks/:id
- [ ] POST /comments (create page-level or rating annotation)
- [ ] GET /comments?page_type=artist&page_id=x (page-level thread)
- [ ] GET /comments?rating_id=x (rating annotation)
- [ ] DELETE /comments/:id (own comments only)
- [ ] POST /folders (create folder)
- [ ] GET /folders?user_id=x (list user's folders)
- [ ] GET /folders/:id (folder contents)
- [ ] PUT /folders/:id (rename, update description, toggle public)
- [ ] DELETE /folders/:id
- [ ] POST /folders/:id/ratings (add rating to folder)
- [ ] DELETE /folders/:id/ratings/:rating_id (remove rating from folder)

**Mobile**
- [ ] Artist detail screen (with page-level comment thread)
- [ ] Album detail screen (with page-level comment thread)
- [ ] Track detail screen (with page-level comment thread)
- [ ] CommentThread component (flat list + reply support for page comments)
- [ ] CommentInput component
- [ ] Rating method picker (comparison vs manual)
- [ ] Beli-style comparison flow (ComparisonCard component)
- [ ] Manual rating screen (star input + decimal + annotation field)
- [ ] Rating annotation visible on profile next to score
- [ ] "Add to folder" action on any rating
- [ ] Folders list screen
- [ ] Folder detail screen (ratings inside folder)
- [ ] Create/edit folder screen
- [ ] Rating history screen (with folder filter)
- [ ] Spotify playlists screen (list all user playlists)
- [ ] Playlist detail screen (track list with rated/unrated status indicator per track)
- [ ] "Rate all unrated" flow — queue up unrated tracks from a playlist for quick sequential rating
- [ ] PlaylistTrackRow component (track name, artist, album art, rating badge if already rated, rate button)

---

### Sprint 4 — Profiles + shareable graphics (Week 7-8)
**Backend**
- [ ] GET /profiles/:username
- [ ] PUT /profiles/:username
- [ ] Cloudflare R2 setup for avatar uploads
- [ ] POST /profiles/avatar

**Mobile**
- [ ] Profile screen (own profile, shows Rank #X + points)
- [ ] Public profile view
- [ ] Profile edit screen
- [ ] Graphics export feature
  - [ ] Top 10 artists card
  - [ ] Monthly recap card
  - [ ] Taste profile card
- [ ] Native share sheet integration (iOS + Android)
- [ ] Artist search screen

**Web (Next.js)**
- [ ] Initialize Next.js on Vercel (rewinds.me)
- [ ] Public profile page (rewinds.me/username)
- [ ] Connect to Supabase for profile + rating data

---

### Sprint 5 — Tour alerts + discovery (Week 9-10)
**Backend**
- [ ] Songkick API integration
- [ ] Ticketmaster API integration
- [ ] POST /tour-alerts (auto-created on rating, manual toggle)
- [ ] GET /tour-alerts?user_id=x
- [ ] Cron job: poll for new tour dates every 6 hours
- [ ] Push notification queue for tour alerts
- [ ] SeatGeek affiliate link generation
- [ ] GET /discover?genre=x (low follower count + rising ELO)

**Mobile**
- [ ] Discover screen (underground rising artists)
- [ ] Tour alert toggle on artist detail screen
- [ ] Tour alert notification handler
- [ ] Tour detail screen (venue, date, city)
- [ ] "Get Tickets" button with SeatGeek affiliate link
- [ ] Expo push notification setup

---

### Sprint 6 — Predictions (Week 11-12)
**Backend**
- [ ] POST /predictions (place bet, deduct wager immediately, compute multiplier)
- [ ] GET /predictions?user_id=x (open + resolved)
- [ ] GET /predictions/available?genre=x (artists with ≥50 votes)
- [ ] Hourly cron: resolve predictions against Redis ranks
- [ ] Credit payout on win + log to point_transactions
- [ ] GET /profiles/:username/stats (win rate, biggest win, total earned)
- [ ] Anti-abuse enforcement

**Mobile**
- [ ] Predictions screen
  - [ ] Open predictions with live rank tracker
  - [ ] Resolved predictions history (won/lost)
  - [ ] New prediction flow: artist → genre → target rank → deadline → wager → payout preview → confirm
- [ ] Rank #X on profile tappable → global user leaderboard
- [ ] Push notifications
  - [ ] Prediction won: "🎯 +X pts"
  - [ ] Prediction lost: "❌ Better luck next time"
  - [ ] "You're X pts from Rank #Y — keep going!"

---

### Sprint 6b — Quests + promo ratings (Week 12)
**Backend**
- [ ] Seed quests table with initial daily + weekly quest definitions
- [ ] GET /quests (all active quests with user's current progress)
- [ ] POST /quests/:id/claim (claim points for completed quest)
- [ ] Quest progress cron: increment user_quest_progress on qualifying actions (votes, ratings, comments)
- [ ] Quest reset cron: create fresh user_quest_progress rows at midnight (daily) and Monday midnight (weekly)
- [ ] GET /promo/available?user_id=x (active campaigns matched to user's genres)
- [ ] POST /promo/:campaign_id/rate (submit promo rating, award pts, enforce one-per-campaign)
- [ ] Admin endpoints for creating/managing campaigns (internal use only for now)

**Mobile**
- [ ] Quest board screen (daily + weekly tabs, progress bars, claim buttons)
- [ ] QuestCard component (title, description, progress, reward, claim CTA)
- [ ] Push notification: "You have uncompleted daily quests" (fires 8pm if quests unfinished)
- [ ] Promo card in discover feed (clearly labeled as promoted)
- [ ] Promo rating flow (artist info + optional preview player + rate + claim points)
- [ ] PromoCard component
- [ ] Points balance updates in real time after quest claim and promo rating

---

### Sprint 7 — Monetization + launch (Week 13-14)
**Backend**
- [ ] Stripe webhook handler
- [ ] POST /subscriptions (Stripe checkout session)
- [ ] RevenueCat entitlement sync
- [ ] Premium middleware (gate features, bump credit/point multipliers)

**Mobile**
- [ ] Premium upgrade screen
- [ ] Stripe web checkout link (US users)
- [ ] StoreKit fallback (non-US)
- [ ] RevenueCat SDK integration
- [ ] Premium features gated: custom profile themes, extra credits, animated graphic exports, higher bet limits
- [ ] Profile theme customization
- [ ] App icon + splash screen
- [ ] Onboarding flow
- [ ] App Store screenshots
- [ ] Privacy policy page (rewinds.me/privacy)

**Launch**
- [ ] Submit to Google Play Store first (faster review)
- [ ] Submit to Apple App Store
- [ ] Post on r/hiphopheads, r/indieheads, r/rnb
- [ ] Set up TikTok, batch first 10 videos
- [ ] Apply to SeatGeek + Ticketmaster affiliate programs