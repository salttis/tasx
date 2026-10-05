# Tasx-projektisuunnitelma

**Päivitetty:** 2026-10-05
**Tila:** Go-versio lukee ja kirjoittaa Markdown-tehtäviä, tarjoaa tekstimuotoisen roadmap-komennon ja kirjoittavan kolmen paneelin TUI:n. v1.0:n versiointi, jakelupaketit ja `tx`-asennuksen palautuspolku on toteutettu. TUI:n vuorovaikutteinen validointi aidossa Windows-terminaalissa on vielä tekemättä ja on julkaisun tarkastusportti.

## Tavoite

Rakentaa Tasxista itsenäinen, local-first Go-tehtävänhallintaohjelma, jossa CLI on ensisijainen rajapinta ja Lazygit-tyylinen TUI nopeuttaa päivittäistä selausta ja tehtävien käsittelyä. Ohjelma tukee ainoastaan Markdown-checkbox-tehtäväformaattia ja toimii paikallisesti ilman palvelua tai tietokantaa.

`tasx` tarjoaa ihmisluettavat komentorivitulosteet. Asennettu `tx`-binäärinimi avaa TUI:n ilman argumentteja.

## Nykytila

Go-repossa on:

- Cobra-pohjainen komentorunko, `list`, `roadmap`, `config`, `version` ja `tui`.
- Markdown-`tasks`-tiedostojen luku ja tekstimuotoinen listaus.
- TOML-muotoiset käyttäjäasetukset tiedostossa `~/.tasxrc`.
- Scopejen valinta käyttäjätason listasta, nykyisestä reposta ja rekisteröidyistä projekteista.
- Bubble Tea- ja Lipgloss-pohjainen kolmen paneelin TUI tehtävien selaamiseen, lisäykseen ja muokkaamiseen; roadmap näkyy yksityiskohtapaneelissa.
- Markdown-checkbox-, osio-, kommentti-, tagi-, projekti- ja päivitysaikamuodon luku.
- Koneellisesti luettavan roadmapin jäsennys ja CLI-tekstiesitys.
- Tehtävien lisäys, valmistuminen, uudelleenavaus, työnkulkutilan muutos sekä käyttäjän käynnistämä arkistointi.
- `tx.exe`-asennusskripti.
- Git-tageihin perustuva v1.0:n monialustainen julkaisuputki, tarkistussummat ja palautettava käyttäjäkohtainen `tx.exe`-asennus.

Roadmap-käskyn lisäksi projektien tilakooste puuttuu. TUI:n interaktiivinen käyttö ei ole vielä vahvistettu oikeassa Windows-terminaalissa; v1.0-julkaisu valmistellaan draftina, kunnes tarkistus on tehty. Yksityiskohtainen rajaus on [tehtäväformaattikuvauksessa](./task-file-format.md).

## Toimintaperiaatteet

1. **Yksi tietolähde per scope:** käyttäjätason `~/.personal/tasks` tai repon `.ai/tasks`. Tiedosto on Markdown-yhteensopiva. Ei tietokantaa, daemonia, `tasks/`-hakemistoa eikä erillistä valmiiden tehtävien arkistoa.
2. **Markdown-only:** Go tukee vain Markdown-checkbox-formaattia.
3. **Ei hiljaista scope-vaihtoa:** eksplisiittinen `--repo` vaatii repo-listan. Automaattinen oletus voi käyttää käyttäjälistausta, jos repo-listaa ei ole.
4. **TUI ei omista domain-logiikkaa:** CLI:n ja TUI:n tulee käyttää samoja tehtävä- ja tallennustoimintoja.
5. **Vain tarpeelliset riippuvuudet:** Cobra CLI:lle, Bubble Tea/Lipgloss TUI:lle ja TOML-konfiguraatiolle. `log/slog` on loggerin API. Shared Go -moduuli on valinnainen eikä Tasxin build saa riippua paikallisesta `replace`-polusta.
6. **Projektikohtaiset ID:t ja säilyttävät kirjoitukset:** tehtävän juokseva numero tallennetaan `$numero`-muodossa ja kvalifioidaan projektilla komennoissa muodossa `projekti-numero`. Kirjoittava versio säilyttää osiot, checkbox-rakenteen, sisennykset, kommentit, rivijärjestyksen, tuntemattomat metatiedot ja muokkaamattomat rivit. Muutetun tehtävän ISO 8601 -UTC-päivitysaika kirjoitetaan `YYYY-MM-DDTHH:MM:SS.sssZ`-muodossa rivin loppuun. Tehtävää ei poisteta tiedostosta; arkistointi on käyttäjän erikseen käynnistämä siirto `Arkisto:`-osioon.

## Onnistumiskriteerit

- Go-versio pystyy lukemaan nykyisiä käyttäjätason ja repo-tason tehtävälistoja muuttamatta niitä.
- CLI:n listaus-, suodatus-, muokkaus-, tilanvaihto-, valmis-/uudelleenavaus-, scope- ja roadmap-käyttötavat ovat dokumentoituja. Tuloste on tekstimuotoinen.
- Kirjoittavat operaatiot muuttavat vain kohderiviä, lisäävät tehtävän tai siirtävät sen arkistoon; virhe ei jätä osittaista tai rikkoutunutta tiedostoa. Tehtävän ID säilyy arkistoinnissa.
- TUI:ssa tehtävät voidaan selata ja käsitellä samoilla sovellusoperaatioilla kuin CLI:ssä. TUI:ta käytetään aidossa terminaalissa ennen sen julistamista valmiiksi.
- `tx`-asennus ja Go-version käyttöönotto eivät muuta PATHia tai korvaa vanhaa `tasx`-komentoa automaattisesti.
- Dokumentaatio kertoo erotetusti, mikä on toteutettu, suunnitteilla ja vielä todentamatta.

## Rajaus

Tässä projektissa ei oteta tavoitteeksi pilvisynkronointia, käyttäjätilejä, monen käyttäjän palvelua, tietokantaa, erillistä done-arkistoa tai selain-/desktop-käyttöliittymää. Tasx käyttää olemassa olevaa tiedosto- ja repo-työnkulkua. Uusi toiminnallisuus lisätään vain todelliseen tarpeeseen.

## Toteutusjärjestys

Vaiheet ja hyväksymiskriteerit ovat [roadmapissa](./roadmap.md). CLI:n kirjoitusoperaatiot, roadmapin tekstirajapinta sekä TUI:n tehtävämuokkaukset on toteutettu. Seuraava tarkistus on TUI:n vuorovaikutteinen käyttö aidossa Windows-terminaalissa.
