# Introduction

A program to generate Anki flashcards for segmented/"phrasal" translations of the Qurʾān.

An example deck can be seen in `out/`. It is based on Saheeh International. Please read the section below on [Accuracy](#accuracy) before using it.

# Installation

```sh
go install git.sr.ht/~rehandaphedar/tarjumlize@latest
```

# Usage

The documentation for usage and flags can be accessed by running `tarjumlize -h`.

- The `-translation` data should be a [rabtize compatible JSON](https://sr.ht/~rehandaphedar/rabtize/#output-format)
- The `-words` data can be obtained from QUL's [Ayah by ayah and word by text of Quran](https://qul.tarteel.ai/resources/quran-script)
- The `-layout` data can be obtained from QUL's [Mushaf Layout Resources](https://qul.tarteel.ai/resources/mushaf-layout)
- The `-metadata-*` can be obtained from QUL's [Quran data, surahs, ayahs, words, juz etc.](https://qul.tarteel.ai/resources/quran-metadata)
- The `-media-config` is a YAML file with a list of objects with the keys `src` and `as`. The filepaths are resolved relative to the config file.

# Card Types

Two card types are generated: `Translation Recall` (Arabic on front; translation on back) and `Arabic Recall` (translation on front; Arabic on back).

# Accuracy

Note that the original purpose of [rabtize](https://sr.ht/~rehandaphedar/rabtize) (and it's sister project, [jumlize](https://sr.ht/~rehandaphedar/jumlize)) was to split Qurʾānic verses and translations into segments so that long verses can fit properly on screens.

Both of the programs use AI.

Thus, it is not recommended to use the output of [rabtize](https://sr.ht/~rehandaphedar/rabtize) directly as an authoritative translation of the *segments* (even if the translation of the verse as a whole might be authoritative).

One should rather use a translation that is segmented/"phrasal" *from the get go*. However, since there is no such translation *in [rabtize](https://sr.ht/~rehandaphedar/rabtize) compatible JSON* yet, the output of rabtize is being used for testing purposes temporarily.

# Arabic Recall Cards Not Showing

Due to a bug in [genanki-go](https://github.com/npcnixel/genanki-go/pull/9), `Arabic Recall` cards are not generated correctly when the deck is imported. As a result, they do not appear in the Browser and are not included in reviews.

A [Pull Request](https://github.com/npcnixel/genanki-go/pull/9) has been opened with the fix.

Until it is merged, run `Tools → Check Database` after importing the deck. This will rebuild the affected cards and make them appear normally.
