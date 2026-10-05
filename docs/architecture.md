# Arkkitehtuuri

## Tavoiterakenne

```text
cmd/tasx
  └── internal/cli ─┐
      internal/tui ─┴── tehtävä- ja scope-operaatiot ── tasks-tiedostot
```

CLI ja TUI ovat esitystapoja. Tehtävän parsinta, tilat, scopevalinta ja tiedostomuutokset kuuluvat niiden ulkopuolelle ja niitä tulee käyttää molemmista samoina.

## Nykyiset paketit

| Polku | Vastuu |
|---|---|
| `cmd/tasx` | Ohjelman entrypoint; binäärin nimi `tx` valitsee oletuksena TUI:n. |
| `internal/cli` | Cobra-komennot: listaus, tehtävien kirjoitus, roadmapin tekstiesitys, arkistointi, TUI:n käynnistys, asetukset ja versio. |
| `internal/config` | `~/.tasxrc` TOML-asetukset ja niiden validointi. |
| `internal/scope` | Nykyisen Git-repon tunnistus, käyttäjätason scope ja rekisteröidyt projektit. |
| `internal/task` | Markdown-checkbox-formaatin luku sekä tehtävien osio-, hierarkia- ja metatietomalli. |
| `internal/roadmap` | Koneellisesti luettavan roadmapin parsinta ja tekstiesityksen kirjoitus. |
| `internal/store` | Markdown-muotoisen `tasks`-tiedoston luku, rivikohtaiset muutokset, lukitus ja atominen tallennus. |
| `internal/tui` | Bubble Tea -tilamalli, kolmen paneelin roadmap-näkymä ja tehtävätoiminnot. |

CLI ja TUI käyttävät samaa tallennuskerrosta tehtävien lisäykseen ja muutoksiin.

## Tallennus ja scopet

- Käyttäjän tehtävälista on `~/.personal/tasks`, ellei käyttäjä muuta juurta asetuksella `personal_dir` tai `TASX_HOME`.
- Repon tehtävälista on `<repo>/.ai/tasks`.
- Rekisteröidyt repot haetaan käyttäjän projektirekisteristä `~/.personal/ai/projects` (tai valitusta personal-hakemistosta).
- Scope ja tila ovat eri suodattimia. `--repo` ja `--global` ovat eksplisiittisiä valintoja.
- Ilman scope-valitsinta nykyinen repo valitaan, jos sen tehtävälista on olemassa; muutoin käytetään käyttäjälistaa. Rekisteröidyn projektin voi valita `list <projekti>`, `add --projekti <projekti>` tai globaalilla `projekti-numero`-tehtävä-ID:llä.
- Go-versio lukee Markdown-checkbox-tehtävät (ID valinnainen). Vanhaa tehtäväformaattia tai muuntoa ei tueta. Ks. [tehtäväformaatti](./task-file-format.md).

## Käyttöliittymä

TUI käyttää Bubble Tea -tapahtumamallia ja Lipgloss-tyylejä. Lazygitin kaltainen käyttömalli jakaa näkymän scope-, tehtävä- ja yksityiskohtapaneeleihin; yksityiskohtapaneeli näyttää valitun tehtävän sekä roadmapin. Kapeassa terminaalissa näytetään valittu paneeli kerrallaan. Tehtävien muutokset suoritetaan tallennuskerroksen kautta. TUI pysyy `internal/tui`-paketin sisällä.

## Riippuvuudet ja lokitus

- Cobra tarjoaa CLI-komennot ja liput.
- Bubble Tea ja Lipgloss tarjoavat TUI:n tapahtumamallin ja esityksen.
- TOML v2 lukee käyttäjäasetukset.
- Standardin `log/slog`-rajapinta on valittu yhteiseksi lokitusrajapinnaksi. Shared-moduulin `logx` ei ole Tasxin pakollinen riippuvuus.
- Yhteiset Go-kirjastot ovat paikallisia, valinnaisia lähdemoduuleja; Tasxin buildin on pysyttävä itsenäisenä.

Lisää uusi yleinen riippuvuus vasta, kun sillä on selkeä, toteutuksessa tarvittava käyttötapaus.
