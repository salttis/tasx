# Tehtävätiedoston formaatti

Tasx käyttää Markdown-checkbox-tehtävämuotoa.

## Markdown-muoto

### Osiot

Otsikoissa voi käyttää Markdownin `#`–`######` tasoja. Myös esimerkin kaltaiset tavalliset otsikot hyväksytään, kun rivi päättyy kaksoispisteeseen; tällaiset otsikot tulkitaan ensimmäisen tason osioiksi. Otsikon kaksoispiste ei kuulu näytettävään osion nimeen.

```markdown
# Kehitys:
## Käyttöliittymä:

### Seuraava julkaisu:
```

Osioiden nimet ja niiden polku liitetään tehtävään. Tyhjä rivi ei päätä osiota. Otsikon jälkeen tulevat tehtävät kuuluvat siihen, kunnes saman tai ylemmän tason otsikko vaihtaa osion.

### Tehtävät ja alitehtävät

Tehtävä on Markdown-checkbox:

```markdown
- [ ] Avoin tehtävä
- [x] Valmis tehtävä
    - [ ] Alitehtävä
        - [ ] Syvempi alitehtävä
```

Sisennys määrittää vanhemman ja syvyyden; alitehtävä voi sisältää omia alitehtäviä. Valinnainen projektikohtainen juokseva ID kirjoitetaan tehtävän tekstin osaksi dollarimerkillä:

```markdown
- [ ] Tämä on tehtävä jossa ID-numero $2
- [ ] Tämä on tehtävä jossa ID-numero $003
- [ ] Tämä on tehtävä jossa ID-numero $22
```

Tunnistettu `$numero` poistetaan näytettävästä kuvauksesta. ID:t ovat scopen sisällä numeerisia; komentorivillä nykyisen projektin tehtävään viitataan pelkällä numerolla (`2`) ja projektin tehtävään mistä tahansa `projekti-2`-muodossa.

### Metatiedot

Metatiedot kirjoitetaan tehtävän kuvauksen perään:

| Muoto | Merkitys |
|---|---|
| `#tag` | Tehtävän tagi. Useita tageja voi käyttää. |
| `@next`, `@waiting`, `@parking`, `@someday`, `@review` | Tehtävän työnkulkutila. Enintään yksi tila kannattaa asettaa. |
| `+projekti` | Projektitagit; useita voi käyttää. |
| `$numero` | Valinnainen projektikohtainen juokseva ID. |
| `YYYY-MM-DDTHH:MM:SS.sssZ` | Valinnainen ISO 8601 -päivitysaika; aikaleima on tehtävärivin viimeinen kenttä. |

Esimerkiksi:

```markdown
- [ ] Lisää tehtävien suodatus #ui #development @next +tasx $666 2026-10-02T00:22:00.000Z
```

Päivitysaika on ISO 8601 / RFC 3339 -aikaleima. Tasx kirjoittaa sen aina UTC-muodossa kolmen millisekuntinumeron tarkkuudella, esimerkiksi `2026-10-05T17:28:00.000Z`. Lukija hyväksyy myös muut RFC 3339 -murto-osat ja aikavyöhykkeet. Aikaleima ei tarkoita luonti- eikä valmistumisaikaa. Valmis-checkbox ilman erillistä valmistumistietoa ei saa keksittyä valmistumispäivää. Aikaleimaksi kelpaamaton tai rivin keskellä oleva teksti säilyy osana kuvausta.

### Kommentit

Sisennetty tavallinen listakohta, joka ei ole checkbox, on edeltävän sisemmän tehtävän kommentti:

```markdown
- [ ] Tallenna suodattimet $666
    - Tämä kommentti kuuluu tehtävälle.
    - Huomio #ui #devops
```

Kommentti on tavallista tekstiä; sen tagit ja projektitagit erotellaan kommentin omiksi metatiedoiksi. Kommentit kuuluvat lähimpään aiempaan tehtävään, jonka sisennys on pienempi kuin kommenttirivin.

## Esimerkki

```markdown
Ominaisuudet:
- [ ] Lisää TUI-haku #ui @next +tasx $100 2026-10-02T00:22:00.000Z
    - Haku kattaa ID:n ja kuvauksen.
    - [ ] Lisää tagisuodatus
        - Valitse useampi tagi #ui #development
```

## Tallennusmalli

Kunkin scopen tehtävien totuuslähde on yksi UTF-8-tekstitiedosto:

- käyttäjäscope: `~\.personal\tasks`
- repo-scope: `<repo>\.ai\tasks`

Go-versio ei käytä tietokantaa. CLI:n tulosteet ovat ihmisluettavaa tekstiä. Kirjoitusoperaatiot säilyttävät tiedoston rivit, rivinvaihdot ja BOM-merkin; ne vaihtavat vain kohderivin kenttiä tai lisäävät/siirtävät tehtävän. Tallennus on atominen, tehtävälistan kirjoittajat sarjallistetaan lukolla ja ennen tallennusta havaittu ulkoinen muutos keskeyttää päivityksen.

## Go-lukijan käyttäytyminen

`internal/task`-parseri:

- tunnistaa Markdown-checkboxit myös ilman ID:tä ja kokoaa osion, syvyyden, vanhemman ID:n, tagit, projektitagit, kommentit sekä viimeisenä olevan UTC-päivitysajan;
- ohittaa otsikot ja kommenttirivit varsinaisina tehtävinä;
- säilyttää jokaisen tehtävän alkuperäisen rivin `Line`-kentässä ja raportoi lähderivin numeron;
- lukee ISO 8601 -päivitysajan rivin viimeisestä kentästä;
- päivittää jokaisen muuttuneen tehtävän rivin lopun aikaleiman muotoon `YYYY-MM-DDTHH:MM:SS.sssZ`.

`tasx add`, `done`, `reopen`, `status` ja `archive` ovat kirjoittavia komentoja. Uusi tehtävä saa juoksevan `$1`, `$2`, … -ID:n kyseisen projektin tehtävälistassa. Nykyisen repon tai käyttäjäscopen tehtäviä voi käsitellä pelkällä numerolla; projektirepon tehtävään viitataan globaalisti `projekti-numero`-muodossa. `tasx add ... --projekti nimi` lisää tehtävän rekisteröityyn projektiin. Valmistuminen ja avaaminen eivät poista tai arkistoi tehtävää.

`tasx archive <id>` siirtää käyttäjän valitseman tehtävän ja sen alitehtävät tiedoston alaosaan `Arkisto:`-osioon. Arkistointi ei ole automaattinen; käyttäjä päättää sen erikseen. Siirretyn tehtävän oma päivitysaika päivittyy, mutta alitehtävien ja kommenttien rivit säilyvät. ID:t eivät muutu, joten ID-sarja jatkuu arkistoitujenkin tehtävien yli. Uudelleenajettu arkistointi ei muuta tiedostoa.

Roadmap luetaan tiedostosta `.ai/roadmap.md` (repo) tai `roadmap.md` (käyttäjäscope). `tasx roadmap [projekti]` näyttää tavoite-, hyväksymiskriteeri-, vaihe- ja tehtäväriippuvuusrakenteen tekstinä.
