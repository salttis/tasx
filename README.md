# Tasx

Tasx on local-first tehtävänhallintaohjelma, joka tukee vain Markdown-checkbox-formaattia. CLI näyttää tulokset tekstinä. Katso [tehtäväformaatti](./docs/task-file-format.md) ja muu [dokumentaatio](./docs/README.md).

## Julkaisu

Tee minor- tai major-julkaisu repojuuresta:

```powershell
.\scripts\release.ps1 minor
.\scripts\release.ps1 major
```

Skripti tarkistaa työpuun ja GitHub-tilan, ajaa tarkistukset, luo version Git-tagin ja julkaisee sen jälkeen GitHub-releasen julkiseksi. Se pyytää kirjoittamaan ehdotetun tagin ennen ulospäin julkaisevia toimia. Tarkemmat vaatimukset ja vaiheet ovat [julkaisuoppaassa](./docs/release.md). Versio tulee Git-tagista; kehitysbuildit näyttävät `dev`.

## Dokumentaatio

Projektin suunnitelma, nykyisen toteutuksen kuvaus ja kehitysohjeet on koottu [docs-hakemistoon](./docs/README.md).

## Kehitysympäristö

Tarvitaan Go 1.27+ ja Git. Paikalliset komennot:

```powershell
go run ./cmd/tasx --help
go run ./cmd/tasx list --repo
go run ./cmd/tasx tui
```

Tasx päättelee oletusscopen ajokontekstista: repossa käytetään repon tehtävälistaa, muualla käyttäjälistaa. `tasx.exe list projektinimi` näyttää rekisteröidyn projektin tehtävät, ja projektikohtaisia tehtäviä voi käsitellä mistä tahansa `projekti-numero`-ID:llä. `tasx.exe roadmap [projektinimi]` näyttää roadmapin tavoite- ja vaiheistuksen tekstinä. `tasx.exe tui` avaa Lazygit-tyylisen kolmen paneelin näppäimistöliittymän tehtävänäkymällä ja roadmap-tiedoilla. `tx` avaa TUI:n ilman argumentteja.

Rakenna ja asenna `tasx`-CLI sekä `tx`-TUI-aliaksen nykyisen käyttäjän henkilökohtaiseen komentohakemistoon:

```powershell
pwsh -NoProfile -File .\scripts\install.ps1
```

Tämä luo `dist\tasx.exe`- ja `dist\tx.exe`-ohjelmat ja kopioi molemmat hakemistoon `~\.personal\scripts`, joka on Tasxin henkilökohtaisten komentojen hakemisto. `tasx.exe` käynnistää CLI:n ja `tx` TUI:n. Asennin ei korvaa mahdollista vanhaa `tasx.ps1`-komentoa; PowerShellissä käytä `tasx.exe`-nimeä, jos vanha skripti on samassa hakemistossa. Ohjelma ei muuta PATH-asetusta. Avaa uusi terminaali tarvittaessa.

Julkaistun version voi asentaa esimerkiksi näin:

```powershell
pwsh -NoProfile -File .\scripts\install.ps1 -Version 1.0.0
```

Jos aiempi `tasx.exe` tai `tx.exe` oli olemassa, asennus säilyttää siitä erillisen palautuskopion. Palautus:

```powershell
pwsh -NoProfile -File .\scripts\install.ps1 -Rollback
```

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
