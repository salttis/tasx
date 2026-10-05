# Kehittäminen

## Vaatimukset

- Go 1.27 tai uudempi.
- Git scope- ja repo-tunnistukseen.
- PowerShell 7 asennusskriptin ajamiseen Windowsissa.

TUI:n paikallinen kehitys ei vaadi .NETiä. Spectre.Console ei kuulu Go-version käyttöliittymäpinoon.

## Tavalliset komennot

Repojuuresta:

```powershell
go build ./...
go run ./cmd/tasx list --repo
go run ./cmd/tasx tui
```

## Visual Studio Code

`.vscode/launch.json` sisältää CLI- ja TUI-debuggausprofiilit. Molemmat käyttävät eksplisiittisesti repon tehtävälistaa. Komennot löytyvät `Terminal > Run Task` -valikosta:

- `Tasx: Run CLI (repo tasks)` ja `Tasx: Run TUI (repo tasks)` käynnistävät sovelluksen.
- `Tasx: Test`, `Tasx: Check formatting`, `Tasx: Vet` ja `Tasx: Build Windows binary` suorittavat vastaavat kehitystarkistukset.
- `Tasx: Release checks` ajaa muotoilu-, testi-, vet- ja build-tarkistukset järjestyksessä.
- `Tasx: Publish release tag (starts draft release)` luo version perusteella annotatoidun tagin ja pushaa sen `origin`-remoteen. Tämä käynnistää GitHub Actionsin draft-julkaisun; suorita tehtävä vain, kun versio on valmis julkaistavaksi.

Asennetun käyttäjäkomennon päivittäminen:

```powershell
pwsh -NoProfile -File .\scripts\install.ps1
```

Julkaisuversion asennus ja edellisen binäärin palautus on kuvattu [julkaisuoppaassa](./release.md).

## Muutosohjeet

- Säilytä CLI ei-interaktiivisena, kun käyttäjä ei käynnistä TUI:ta.
- Pidä TUI-framework käyttöliittymäpaketissa ja jaa sovellusoperaatiot CLI:n kanssa.
- Älä muuta tehtävätiedostoa suoraan osana dokumentaatiota tai listauskomentoja.
- Käsittele vanhoja rivejä säilyttävästi; älä pudota tuntemattomia kenttiä.
- Tasx tukee vain Markdown-checkbox-tehtäväformaattia; älä lisää vanhan formaatin parseria tai migraatiopolkuja.
- Älä lisää riippuvuutta, jos nykyisen kirjaston tai standardikirjaston ratkaisu riittää.
- Päivitä tätä dokumentaatiota, projektisuunnitelmaa tai roadmapia, kun toteutettu käyttäytyminen tai vaihe muuttuu.

## Validointi

Valitse muutoksen pääpolkua vastaava rajattu tarkistus. Dokumentaatiomuutokset eivät vaadi testisarjaa. CLI:n tai TUI:n toiminnan muuttuessa build yksin ei todenna käyttäjäkokemusta tai tiedostoyhteensopivuutta; tarkista myös muokattu toiminto sopivassa oikeassa ympäristössä. Interaktiivisen TUI:n käyttöä ei pidä väittää todennetuksi ilman oikeaa terminaalia.
