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

Asennetun käyttäjäkomennon päivittäminen:

```powershell
pwsh -NoProfile -File .\scripts\install.ps1
```

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
