# Tasx-roadmap

Roadmap kuvaa tavoitellut tulokset; se ei muuta tehtävien tiloja. Koneellisesti luettava vaiheistus on [.ai/roadmap.md](../.ai/roadmap.md)-tiedostossa, jonka `tasx roadmap` esittää tekstinä.

## Päämäärä

Tasxista tulee luotettava, itsenäinen Go CLI ja näppäimistökäyttöinen TUI Markdown-tehtäväformaattiin. Vanhaa tehtäväformaattia tai migraatiota ei tueta.

## Vaihe 0 — Perusta ja projektin dokumentaatio

**Tulos:** erillinen Go-repo, nykytilaa kuvaava dokumentaatio, koneellisesti luettava roadmap ja CLI/TUI.

**Tila:** dokumentaatiopaketti, roadmapin Go-parseri ja tekstikomento, v1.0-julkaisupaketit sekä palautettava Windows-asennus on toteutettu. TUI:n interaktiivista toimintaa ei ole vielä todennettu aidossa Windows-terminaalissa.

**Hyväksymiskriteerit:**

- [x] Dokumentaatio kuvaa Go-version käyttörajat.
- [x] Projektisuunnitelma ja vaiheistettu roadmap ovat saatavilla.
- [x] Koneellisesti luettavat roadmap-vaiheet ovat saatavilla `.ai/roadmap.md`-tiedostossa.

## Vaihe 1 — Markdown-tehtäväformaatti

**Tulos:** Tasx tukee vain Markdown-checkbox-tehtäväformaattia. Vanhan formaatin tuki ja muunto poistetaan.

**Hyväksymiskriteerit:**

- [x] Markdown-tehtävät, metatiedot, osiot ja alitehtävät luetaan ja kirjoitetaan.
- [x] Vanhan formaatin parseri ja migraatiokomento eivät kuulu ohjelmaan.

## Vaihe 1a — Tehtävämetatiedot ja aikaleimat

**Tulos:** Go lukee Markdown-checkboxit, osiot, alitehtävät, kommentit, tagit, projektitunnisteet ja rivin lopussa olevan päivitysajan.

**Hyväksymiskriteerit:**

- [x] ID:tön Markdown-checkbox näkyy tehtävänä.
- [x] Osio, sisäkkäisyys, kommentit, `#tag`, `@tila`, `+projekti` ja ISO 8601 -aika tulkitaan erillisiksi kentiksi.
- [x] Uuden formaatin lukukäyttäytyminen tarkistetaan käyttäjän esimerkillä live-ajossa.
- [x] Kirjoittavat komennot päivittävät tehtävärivin lopussa olevan ISO 8601 -UTC-aikaleiman.

## Vaihe 2 — Säilyttävät kirjoitusoperaatiot ja CLI-vastaavuus

**Tulos:** Go-versio voi tehdä tehtävien arkipäiväiset muutokset luotettavasti ja vastaa vanhan CLI:n ydinkomentoja.

**Hyväksymiskriteerit:**

- [x] Lisäys, näyttö, valmis-/uudelleenavaus ja työnkulkutilan vaihto on toteutettu.
- [x] Checkbox-valmius ja päivitysaika säilyvät tehtävärivillä; päivitysaika muuttuu tehtävää muokattaessa.
- [x] Atominen kirjoitus ja samanaikaisten muutosten suojaus on toteutettu ennen kirjoittavan CLI:n käyttöönottoa.
- [x] Kommentit, muokkaamattomat rivit, tuntemattomat kentät ja alkuperäinen järjestys säilyvät.
- [x] Virhetilanteet raportoidaan eikä onnistumisen näköistä osittaista tallennusta jää.
- [x] Arkistointi on vain käyttäjän erikseen käynnistämä toiminto; tehtävien ID:t säilyvät ja juokseva numerointi jatkuu arkistosta.

## Vaihe 3 — Roadmap-rajapinta ja scope-näkymät

**Tulos:** Go lukee koneellisesti roadmapin tavoitteet, hyväksymiskriteerit, vaiheet ja tehtäväriippuvuudet sekä näyttää ne tekstinä. Projektien yhteinen tehtäväkooste ei kuulu tähän toteutukseen.

**Hyväksymiskriteerit:**

- [x] `tasx roadmap [projekti]` näyttää goal-, acceptance criteria-, phase- ja task-riippuvuustiedot.
- [x] Roadmapin lukuvirheet ja puuttuvat tiedostot näytetään virheinä.
- [x] CLI:n tuloste on ihmisluettava teksti; JSON-ulostuloa ei ole.
- [x] Tehtävän tila- ja scope-valinta pysyvät erillisinä.

## Vaihe 4 — Kirjoittava Lazygit-tyylinen TUI

**Tulos:** TUI tarjoaa tehtävien käsittelyn näkyvien scope-, tehtävä- ja yksityiskohtapaneelien kautta.

**Hyväksymiskriteerit:**

- [x] TUI käyttää samaa tallennuskerrosta kuin CLI eikä kirjoita omaa erillistä tallennuslogiikkaa.
- [x] Tehtävän lisäys, valmistuminen, uudelleenavaus ja tilanvaihdot ovat käytettävissä näppäimistöltä.
- [x] Haku, suodatus, scopen vaihto, tyhjät/puuttuvat scopet, virheet ja kapeat terminaalikoot on toteutettu.
- [ ] Interaktiivinen käyttö tarkistetaan Windowsin aidossa terminaalissa ennen v1.0-julkaisun julkaisemista.
- [x] `tx`-käynnistys dokumentoidaan eikä asennus vaihda PATHia hiljaisesti.

## Vaihe 5 — Käyttäjäkomennon käyttöönotto ja vakautus

**Tulos:** Go-versio on riittävän yhteensopiva, että käyttäjän komentopolku voidaan siirtää hallitusti.

**Hyväksymiskriteerit:**

- [x] CLI:n JSON-ulostuloa ei toteuteta; tekstikäyttöliittymä on dokumentoitu.
- [x] `tx`-asennus on versionoitavissa, toistettava ja palautettavissa edelliseen binääriin.
- [x] Tagipohjainen julkaisuputki rakentaa monialustaiset binäärit, tarkistussummat ja tarkastettavan draft-julkaisun.
- [x] Jakelupaketti sisältää erilliset `tasx`- ja `tx`-binäärinimet; `tasx`-oletuslistaus on vahvistettu CLI-ajossa. `tx`-TUI:n interaktiivinen käynnistys on vielä vaiheessa 4 lueteltu julkaisuportti.
- [x] Vanha PowerShell-toteutus säilytetään, kunnes käyttäjä erikseen ottaa Go-komentopolun käyttöön.

## Aikataulu ja priorisointi

Roadmapissa ei luvata kalenteripäiviä eikä tehdä automaattista prioriteettipäätöstä myöhemmistä vaiheista. Projektin roadmap kuvaa tavoitellut tulokset ja toteutustilan; se ei ylläpidä erillistä tehtäväjonoa.
