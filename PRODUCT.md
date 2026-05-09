# Rewinds — Product & Business

## What Rewinds is

A music taste platform targeting underground music fans who care deeply about discovering artists early. The core identity is a **music taste passport** — your ratings, predictions, rank, and discoveries all in one shareable profile.

The closest comparisons are Letterboxd (logged taste + social identity) and fantasy sports (prediction/betting mechanic). No direct competitor does both.

**Domain:** rewinds.me

---

## Why each feature exists

**ELO head-to-head voting**
People find it hard to rate music with a number but easy to say "I like A more than B." ELO voting generates the entire artist ranking system from community preferences rather than editorial decisions. It's also inherently gamified — every vote feels like a decision that matters.

**Beli-style comparison rating**
Same insight applied to personal ratings. Instead of picking 7.4/10, you compare against things you've already rated. Keeps your ratings consistent over time and feels more like a game than a form. Manual rating always available as an alternative for users who prefer it.

**Prediction/betting system**
The stickiest retention mechanic. Users bet points that an artist will climb the leaderboard, with payouts scaled by how far the artist has to climb. Creates a reason to check the app daily to track open predictions. Leaderboard rank based on current points balance (not lifetime) means betting is always a real risk — you can fall in rank.

**Shareable graphics**
Every time someone posts their top 10 card or taste profile on social media, it's a free ad. Rate Your Music and AOTY have zero native sharing. Letterboxd added this and it grew massively on Twitter/X. Our target demographic (underground music fans) is extremely online and loves showing off taste.

**Customizable public profiles**
People link VSCOs in their Instagram bios not because VSCO is better than Instagram but because it expresses a different, more curated side of them. rewinds.me/username serves the same purpose for music taste. Every profile is passive distribution.

**Tour alerts**
Best monetization feature. Affiliate revenue happens automatically when users do something they already wanted to do (buy concert tickets). No friction, no convincing anyone to pay. SeatGeek pays ~$5 flat per ticket — average ticket is $100+.

**Underground artist focus**
AOTY and Rate Your Music treat underground artists as an afterthought. Our ELO system naturally surfaces rising underground artists through the "rising movers" feed. The tastemaker identity ("I was listening to them before they blew up") is core to our target user's self-image.

---

## Target user

Underground music fans aged 18–28 who:
- Follow music on Twitter/X, Reddit (r/hiphopheads, r/indieheads), TikTok
- Pride themselves on discovering artists early
- Are frustrated that existing tools (AOTY, Rate Your Music) have bad UX and no mobile-first experience
- Already share their music taste on social media but have no dedicated platform for it

---

**Quest system**
Quests replace passive per-action point rewards. Instead of earning points for every vote or rating (which creates mindless farming), users complete intentional daily and weekly challenges. This produces higher quality engagement — someone completing a "rate 5 albums in the same genre" quest is more engaged than someone tapping through 50 votes to farm 500 points. The quest board also creates a natural daily re-engagement hook and a powerful push notification reason ("You have 2 daily quests remaining").

**Artist promo campaigns**
Labels and artists pay to surface new releases to taste-matched users. Users earn their highest single point reward by completing a promo rating — rate the promoted artist (preview optional). This creates a clean three-way value exchange: label gets confirmed ratings from relevant listeners, user gets points, Rewinds gets revenue. One claim per user per campaign prevents gaming.

---

## Monetization

| Stream | How | Timeline |
|---|---|---|
| Premium subscription | $5/month via Stripe web checkout (US, zero Apple commission) + StoreKit (non-US). Gates higher bet limits, premium quest rewards (1.5x multiplier), profile themes, animated graphic exports. Managed via RevenueCat. | Sprint 7 |
| Tour alert affiliates | SeatGeek affiliate (~$5/ticket flat fee) + Ticketmaster (2–4%) on ticket purchases via push notification links. Completely passive — revenue happens when users buy tickets they already wanted. | Sprint 5 |
| Artist/label promo campaigns | Labels pay per confirmed rating from taste-matched users. Pricing model: guaranteed ratings (e.g. $200 for 100 confirmed ratings from Hip Hop listeners). Campaign ends when rating quota is hit. Requires sales outreach — viable after 5k active users. | Post-launch |
| Affiliate links | Spotify, Apple Music, Discogs vinyl links on every artist and album page. Passive, small amounts at small scale. | V2 |

### Promo campaign pricing model
Generic "boost to music fans" impressions → worth $50.
"100 confirmed ratings from Rewinds Hip Hop listeners with taste profiles matched to your artist" → worth $200–500.

The targeting and confirmed rating requirement is the product. Charge per rating delivered, not per impression. Labels only pay when a user actually engages.

**Suggested pricing tiers (post-launch):**
| Package | Ratings guaranteed | Genre targeting | Price |
|---|---|---|---|
| Starter | 50 ratings | Yes | $99 |
| Standard | 200 ratings | Yes | $299 |
| Premium | 500 ratings | Yes + featured placement | $599 |


### Stripe on iOS
As of the April 2025 Epic v. Apple ruling, U.S. App Store apps can link to an external Stripe web checkout with zero Apple commission. This means:
- US users go through Stripe web checkout → you keep ~97% of revenue (Stripe takes 2.9% + $0.30)
- Non-US users use Apple's StoreKit (15–30% commission)
- RevenueCat manages entitlements across both paths so premium features unlock correctly regardless of payment method
- At $5/month: Stripe = $4.85 kept vs StoreKit = $3.50 kept

---

## Marketing

| Channel | Tactic |
|---|---|
| TikTok | AI-generated "unusual person reacts to underground artist" videos (old man, Buddhist monk, etc.). Batch produce, post daily. Low cost, high shareability. |
| Reddit | Post weekly genre leaderboard results to r/hiphopheads, r/indieheads, r/rnb. Genuine value, not spam. |
| Twitter/X | Tastemaker culture is native to music Twitter. Post leaderboard updates, rising artists, "artists to watch" threads. |
| Organic | Shareable profile graphics and rating cards posted to social media drive installs passively. Make the graphics look exceptional — this is the most important growth lever. |
| Cold outreach | Underground music newsletters and small music blogs. Offer to be featured as a discovery tool. |

### Distribution advantage
Every public profile (rewinds.me/username) shared in a bio or post is passive distribution. The more the profiles look good, the more people click and sign up. Profile aesthetic is a product priority, not just a nice-to-have.

---

## Competitive landscape

| App | Weakness | Our angle |
|---|---|---|
| Rate Your Music | Ancient UX, desktop-only, no mobile app | Mobile-first, modern UX |
| Album of the Year | No gamification, no social, no sharing | Predictions, points, shareable graphics |
| Podiums | Generalist (movies + music + TV), no predictions, ~10k downloads | Music-only, underground focus, prediction system |
| Letterboxd | Film only | Same model applied to music |
| Last.fm | Abandoned, passive scrobbling only | Active curation and community ranking |

---

## Points economy rationale

**Why current balance instead of lifetime earned:**
Lifetime points means rank only goes up — no risk to betting. Current balance means every bet is a real decision. Someone who bets recklessly can fall in rank. This makes the leaderboard dynamic and gives veteran users something to defend.

**Why 1,000 signup points:**
Enough to make a few predictions immediately and get invested in the outcome. Not enough to meaningfully affect the leaderboard. Gets new users into the prediction loop before they've earned points organically.

**Why premium gets a 1.5x prediction multiplier:**
Doesn't feel pay-to-win because points have no real monetary value. But it's a meaningful perk for engaged users who bet frequently — the best tastemakers are premium users who are also the most active predictors.

---

## Future features (post-launch)

- **Listening eras** — organize your taste by phase of life ("high school era", "2021 hyperpop phase")
- **Taste compatibility score** — how similar is your taste to another user's, based on ratings and votes
- **Annual taste recap** — Spotify Wrapped-style yearly summary, shareable, drives organic social posts in December
- **Genre-specific prediction leaderboards** — top predictors in Hip Hop vs Indie separately
- **"Found them first" trophy** — permanent badge when an artist you predicted early crosses a follower milestone