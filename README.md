# Tasx

Tasx on local-first tehtävänhallintaohjelma, joka tukee vain Markdown-checkbox-formaattia. Vanhaa formaattia tai migraatiota ei tueta; CLI näyttää tulokset tekstinä eikä tuota JSON-ulostuloja. Katso [tehtäväformaatti](./docs/task-file-format.md) ja muu [dokumentaatio](./docs/README.md).

## Dokumentaatio

Projektin suunnitelma, nykyisen toteutuksen kuvaus ja kehitysohjeet on koottu [docs-hakemistoon](./docs/README.md).

## Kehitysympäristö

Tarvitaan Go 1.27+ ja Git. Paikalliset komennot:

```powershell
go run ./cmd/tasx --help
go run ./cmd/tasx list --repo
go run ./cmd/tasx tui
```

Tasx päättelee oletusscopen ajokontekstista: repossa käytetään repon tehtävälistaa, muualla käyttäjälistaa. `tasx list projektinimi` näyttää rekisteröidyn projektin tehtävät, ja projektikohtaisia tehtäviä voi käsitellä mistä tahansa `projekti-numero`-ID:llä. `tasx roadmap [projektinimi]` näyttää roadmapin tavoite- ja vaiheistuksen tekstinä. `tasx tui` avaa Lazygit-tyylisen kolmen paneelin näppäimistöliittymän tehtävänäkymällä ja roadmap-tiedoilla. `tx` avaa TUI:n ilman argumentteja.

Rakenna ja asenna `tx` nykyisen käyttäjän henkilökohtaiseen komentohakemistoon:

```powershell
pwsh -NoProfile -File .\scripts\install.ps1
```

Tämä luo `dist\tx.exe`-ohjelman ja kopioi sen hakemistoon `~\.personal\scripts`, joka on Tasxin henkilökohtaisten komentojen hakemisto. Ohjelma ei muuta PATH-asetusta. Avaa uusi terminaali tarvittaessa.

## Käyttäjäkohtaiset oletukset

Valinnainen TOML-asetustiedosto on `~\.tasxrc`. Puuttuva tiedosto tarkoittaa sisäänrakennettuja oletuksia; virheellinen asetus tai tuntematon avain ilmoitetaan virheenä.

```toml
personal_dir = "" # tyhjä = ~/.personal; TASX_HOME ohittaa tämän
default_scope = "auto" # auto, repo, global
state = "all" # all, open, done
```

`TASX_CONFIG` vaihtaa asetustiedoston polun ja `TASX_HOME` käyttäjätason Tasx-datan juuren. CLI-valitsimet ohittavat asetustiedoston oletukset. `default_scope = "auto"` valitsee nykyisen Git-repon, kun siellä on `.ai\tasks`; muuten käytetään käyttäjälistausta. `--repo` ei koskaan putoa hiljaa käyttäjäscopeen.

## TUI

Käyttöliittymässä on scope-, tehtävä- ja yksityiskohtapaneelit. `Tab` vaihtaa paneelia, `j/k` tai nuolet liikuttavat valintaa, `/` aloittaa haun, `a` lisää tehtävän, `Space` vaihtaa valmiuden, `s` kierrättää työnkulkutilaa, `f` vaihtaa suodatinta, `r` päivittää ja `q` sulkee. Kapeassa terminaalissa paneelit näytetään vuorotellen. TUI:n käyttöä aidossa Windows-terminaalissa ei ole vielä todennettu.

CLI-komennot `add`, `done`, `reopen` ja `status` sekä TUI käyttävät samaa tallennuskerrosta. Uuden tehtävän ID on rivin lopun `$numero`, juokseva projektikohtainen numero. `archive` siirtää tehtävän alitehtävineen käyttäjän päätöksellä `Arkisto:`-osioon; tehtäviä ei poisteta valmistumisen yhteydessä.

## Suunnitteluperiaatteet

- Yksi ihmisluettava `tasks`-tiedosto per scope; ei tietokantaa, daemonia tai erillistä valmistuneiden arkistoa.
- CLI ja TUI jakavat saman sovelluslogiikan.
- TUI-framework pysyy käyttöliittymäpaketissa; valittu framework on Bubble Tea.
- Yhteiset Go-kirjastot ovat valinnaisia. Ulkoista yhteistä lähdemoduulia ei vaadita Tasxin buildissa.
- Loggerina käytetään Go:n `log/slog`-rajapintaa. Jaettu `shared/go/logx` lisää yhteiset tulostus- ja vaihekäytännöt; Tasx ei riipu konekohtaisesta shared-moduulista.
