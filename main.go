package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"git.sr.ht/~rehandaphedar/genanki-go-utils/v2/pkg/qul"
	"github.com/npcnixel/genanki-go"
	"go.yaml.in/yaml/v4"
)

func main() {
	modelId := flag.Int64("model-id", int64(1421329585), "ID of the model")
	modelName := flag.String("model-name", "tarjumlize", "Name of the model")

	deckId := flag.Int64("deck-id", int64(1823518270), "ID of the deck")
	deckName := flag.String("deck-name", "tarjumlize", "Name of the peck")
	deckDescription := flag.String("deck-description", "Recall the Arabic and the translation of the segment.", "Description of the deck")

	outputPath := flag.String("output", "out/tarjumlize.apkg", "Output filepath")

	templateHtmlPath := flag.String("template-html", "templates/index.gohtml", "Path to template file")
	templateCssPath := flag.String("template-css", "templates/style.css", "Path to CSS file")

	templateQfmtTranslationPath := flag.String("template-qfmt-translation", "templates/qfmt_translation.html", "Path to template Qfmt file for Translation Recall")
	templateAfmtTranslationPath := flag.String("template-afmt-translation", "templates/afmt_translation.html", "Path to template Afmt file for Translation Recall")
	templateQfmtArabicPath := flag.String("template-qfmt-arabic", "templates/qfmt_arabic.html", "Path to template Qfmt file for Arabic Recall")
	templateAfmtArabicPath := flag.String("template-afmt-arabic", "templates/afmt_arabic.html", "Path to template Afmt file for Arabic Recall")

	templateFrontTranslationName := flag.String("template-front-translation", "front_translation", "Name of the front template for Translation Recall")
	templateBackTranslationName := flag.String("template-back-translation", "back_translation", "Name of the back template for Translation Recall")
	templateFrontArabicName := flag.String("template-front-arabic", "front_arabic", "Name of the front template for Arabic Recall")
	templateBackArabicName := flag.String("template-back-arabic", "back_arabic", "Name of the back template for Arabic Recall")

	wordsPath := flag.String("words", "data/qpc-hafs-word-by-word.json", "Path to words data")
	translationPath := flag.String("translation", "data/en-sahih-international-simple.json", "Path to translation data")
	layoutPath := flag.String("layout", "data/qpc-v4-tajweed-15-lines.db", "Path to layout data")
	metadataAyahPath := flag.String("metadata-ayah", "data/quran-metadata-ayah.json", "Path to ayah metadata")
	metadataJuzPath := flag.String("metadata-juz", "data/quran-metadata-juz.json", "Path to juz metadata")
	metadataHizbPath := flag.String("metadata-hizb", "data/quran-metadata-hizb.json", "Path to hizb metadata")
	metadataRubPath := flag.String("metadata-rub", "data/quran-metadata-rub.json", "Path to rub metadata")
	metadataManzilPath := flag.String("metadata-manzil", "data/quran-metadata-manzil.json", "Path to manzil metadata")
	metadataRukuPath := flag.String("metadata-ruku", "data/quran-metadata-ruku.json", "Path to ruku metadata")

	var tagFormat qul.TagFormat

	tagFormat.Chapter = flag.String("tag-format-chapter", "quran::chapter::%03d", "Format of the chapter tag. %d is replaced with the chapter number.")
	tagFormat.Verse = flag.String("tag-format-verse", "quran::verse::%s", "Format of the verse tag. %s is replaced with the zero padded verse key (Example: 001:001).")
	tagFormat.Page = flag.String("tag-format-page", "quran::page::%03d", "Format of the page tag. %d is replaced with the page number.")
	tagFormat.Juz = flag.String("tag-format-juz", "quran::juz::%02d", "Format of the juz tag. %d is replaced with the juz number.")
	tagFormat.Hizb = flag.String("tag-format-hizb", "quran::hizb::%02d", "Format of the hizb tag. %d is replaced with the hizb number.")
	tagFormat.Rub = flag.String("tag-format-rub", "quran::rub::%03d", "Format of the rub tag. %d is replaced with the rub number.")
	tagFormat.Manzil = flag.String("tag-format-manzil", "quran::manzil::%d", "Format of the manzil tag. %d is replaced with the manzil number.")
	tagFormat.Ruku = flag.String("tag-format-ruku", "quran::ruku::%03d", "Format of the ruku tag. %d is replaced with the ruku number.")

	mediaConfigPath := flag.String("media-config", "media/config.yaml", "Path to media config")

	flag.Parse()

	var words map[string]qul.Word
	var translation map[string]Verse
	var metadataAyah map[string]qul.MetadataAyah

	var metadataDivision qul.MetadataDivision

	err := loadJSON(*wordsPath, &words)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*translationPath, &translation)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataAyahPath, &metadataAyah)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataJuzPath, &metadataDivision.Juz)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataHizbPath, &metadataDivision.Hizb)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataRubPath, &metadataDivision.Rub)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataManzilPath, &metadataDivision.Manzil)
	if err != nil {
		log.Fatal(err)
	}
	err = loadJSON(*metadataRukuPath, &metadataDivision.Ruku)
	if err != nil {
		log.Fatal(err)
	}

	index, err := qul.BuildIndex(*layoutPath, words, metadataDivision, tagFormat)
	if err != nil {
		log.Fatalf("build index: %v", err)
	}

	metadataAyahByVerseKey := make(map[string]qul.MetadataAyah)
	for _, metadataAyahEntry := range metadataAyah {
		metadataAyahByVerseKey[metadataAyahEntry.VerseKey] = metadataAyahEntry
	}

	qfmtTranslation, err := readFile(*templateQfmtTranslationPath)
	if err != nil {
		log.Fatal(err)
	}
	afmtTranslation, err := readFile(*templateAfmtTranslationPath)
	if err != nil {
		log.Fatal(err)
	}
	qfmtArabic, err := readFile(*templateQfmtArabicPath)
	if err != nil {
		log.Fatal(err)
	}
	afmtArabic, err := readFile(*templateAfmtArabicPath)
	if err != nil {
		log.Fatal(err)
	}

	css, err := readFile(*templateCssPath)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles(*templateHtmlPath)
	if err != nil {
		log.Fatalf("parse template files: %v", err)
	}
	var buf bytes.Buffer

	model := genanki.NewModel(*modelId, *modelName).
		SetCSS(css).
		AddField(genanki.Field{Name: "SegmentID"}).
		AddField(genanki.Field{Name: "FrontTranslation"}).
		AddField(genanki.Field{Name: "BackTranslation"}).
		AddField(genanki.Field{Name: "FrontArabic"}).
		AddField(genanki.Field{Name: "BackArabic"}).
		AddField(genanki.Field{Name: "Notes"}).
		AddTemplate(genanki.Template{
			Name: "Translation Recall",
			Qfmt: qfmtTranslation,
			Afmt: afmtTranslation,
		}).
		AddTemplate(genanki.Template{
			Name: "Arabic Recall",
			Qfmt: qfmtArabic,
			Afmt: afmtArabic,
		})
	deck := genanki.NewDeck(*deckId, *deckName, *deckDescription)

	for verseKey, verse := range translation {
		// TODO: Similar Verses?
		// instances := renderInstances(index.Word, metadataAyahByVerseKey, phrase)

		templateData := TemplateData{
			VerseKey: verseKey,
		}

		for _, segment := range verse.Segments {

			templateData.Arabic = renderSegment(index.Word, verseKey, segment.WordRange.Start, segment.WordRange.End)
			templateData.Translation = segment.Text

			templateErrorMessage := "error while executing template %s with data %+v: %v"

			err := tmpl.ExecuteTemplate(&buf, *templateFrontTranslationName, templateData)
			if err != nil {
				log.Printf(templateErrorMessage, *templateFrontTranslationName, templateData, err)
			}
			frontTranslation := buf.String()
			buf.Reset()

			err = tmpl.ExecuteTemplate(&buf, *templateBackTranslationName, templateData)
			if err != nil {
				log.Printf(templateErrorMessage, *templateBackTranslationName, templateData, err)
			}
			backTranslation := buf.String()
			buf.Reset()

			err = tmpl.ExecuteTemplate(&buf, *templateFrontArabicName, templateData)
			if err != nil {
				log.Printf(templateErrorMessage, *templateFrontArabicName, templateData, err)
			}
			frontArabic := buf.String()
			buf.Reset()

			err = tmpl.ExecuteTemplate(&buf, *templateBackArabicName, templateData)
			if err != nil {
				log.Printf(templateErrorMessage, *templateBackArabicName, templateData, err)
			}
			backArabic := buf.String()
			buf.Reset()

			segmentId := fmt.Sprintf("%s_%d_%d", verseKey, segment.WordRange.Start, segment.WordRange.End)

			note := genanki.NewNote(
				model.ID,
				[]string{
					segmentId,
					frontTranslation,
					backTranslation,
					frontArabic,
					backArabic,
					"",
				},
				index.Tag.Verse[verseKey],
			)

			noteIdBase := fmt.Sprintf("%d_%s", model.ID, segmentId)
			note.ID = qul.GenerateID(noteIdBase)
			deck.AddNote(note)
		}
	}

	pkg := genanki.NewPackage([]*genanki.Deck{deck}).AddModel(model)

	if *mediaConfigPath != "" {
		mediaConfigDir := filepath.Dir(*mediaConfigPath)

		mediaConfigData, err := os.ReadFile(*mediaConfigPath)
		if err != nil {
			log.Fatalf("read media config: %v", err)
		}

		var mediaEntries []MediaEntry
		if err := yaml.Unmarshal(mediaConfigData, &mediaEntries); err != nil {
			log.Fatalf("parse media config: %v", err)
		}

		for _, mediaEntry := range mediaEntries {
			src := filepath.Join(mediaConfigDir, mediaEntry.Src)
			as := mediaEntry.As
			mediaEntryData, err := os.ReadFile(src)
			if err != nil {
				log.Fatalf("read media entry %s: %v", src, err)
			}
			pkg.AddMedia(as, mediaEntryData)
		}
	}

	if err := pkg.WriteToFile(*outputPath); err != nil {
		log.Fatalf("write package to %s: %v", *outputPath, err)
	}
}
