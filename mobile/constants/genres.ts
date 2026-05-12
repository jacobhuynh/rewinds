export const GENRES = [
  "hiphop",
  "rnb",
  "pop",
  "rock",
  "indie",
  "electronic",
  "jazz",
  "classical",
  "country",
  "latin",
  "metal",
  "folk",
] as const;

export type Genre = (typeof GENRES)[number];
