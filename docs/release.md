# Julkaisu

## v1.0.0 sisältö

v1.0.0 on paikallinen Go-sovellus, joka tarjoaa Markdown-checkbox-tehtävien CLI-hallinnan ja Bubble Tea + Lipgloss -TUI:n. Julkaisu sisältää Windows-, Linux- ja macOS-binäärit `amd64`- ja `arm64`-arkkitehtuureille, `tasx`- ja `tx`-nimiset käynnistimet, sovellus- ja kolmannen osapuolen lisenssitiedot sekä SHA-256-tarkistussummat. Tasxin oma lisenssi on MIT.

Vanhaa tehtäväformaattia, migraatiota ja JSON-ulostuloa ei ole. Vanhaa PowerShell-komentopolkua ei korvata automaattisesti.

## Julkaisuportit

Ennen draft-julkaisun hyväksymistä:

1. Varmista, että työpuu sisältää vain aiottua julkaistavaa muutosta.
2. Aja CI:n Go-muotoilu-, `go vet`- ja build-tarkistukset.
3. Tarkista `tasx version`, `tasx`-oletuslistaus sekä `tx`-käynnistyksen TUI oikeassa Windows-terminaalissa. Kokeile vähintään scopen ja tehtävän selailua, hakua, lisäystä, valmiuden vaihtoa, tilan vaihtoa, päivitystä ja poistumista eristetyllä tehtävälistalla.
4. Tarkista draft-julkaisun paketit, mukana tulevat lisenssitiedostot ja SHA-256-tarkistussummat ennen kuin julkaiset sen GitHubissa.

TUI:n interaktiivista Windows-validointia ei ole vielä tehty tässä kehitysympäristössä. Build tai CLI-ajot eivät korvaa tätä tarkistusta.

## Versioitu paikallisasennus

PowerShell 7 -repojuuresta:

```powershell
pwsh -NoProfile -File .\scripts\install.ps1 -Version 1.0.0
```

Asennus rakentaa versionumeron binääriin, luo `dist\tx.exe`-artefaktin ja kopioi sen `~\.personal\scripts\tx.exe`-polkuun. Se ei muuta PATHia. Jos kohteessa oli jo `tx.exe`, edellinen binääri tallennetaan nimellä `tx.exe.previous`; jos tiedosto oli ennestään olemassa, uusi varmuuskopio saa aikaleimatun nimen. Palautus vaihtaa nykyisen ja edellisen binäärin keskenään:

```powershell
pwsh -NoProfile -File .\scripts\install.ps1 -Rollback
```

Palautus ei poista asennusta vaan vaihtaa kahden viimeisimmän binäärin paikkaa, joten komento voidaan tarvittaessa suorittaa uudelleen.

Kehitysbuildin asentaminen ilman versiota näyttää versiona `dev`. Julkaisuversion voi antaa myös `v1.0.0`-muodossa; asennin tallentaa binääriin `1.0.0`.

## Ladatun paketin käyttö

Paketin voi purkaa Windowsissa PowerShellillä:

```powershell
tar -xzf .\tasx_1.0.0_windows_amd64.tar.gz
.\tasx_1.0.0_windows_amd64\tasx.exe version
.\tasx_1.0.0_windows_amd64\tx.exe
```

Linuxissa tai macOS:ssä:

```sh
tar -xzf tasx_1.0.0_linux_amd64.tar.gz
./tasx_1.0.0_linux_amd64/tasx version
./tasx_1.0.0_linux_amd64/tx
```

Valitse paketin nimi oman käyttöjärjestelmän ja prosessoriarkkitehtuurin mukaan. `third_party_licenses/` sisältää linkitettyjen riippuvuuksien lisenssidokumentit.

## GitHub-julkaisun valmistelu

CI suorittaa pushissa, pull requestissa ja manuaalisesti muotoilu-, `go vet`- ja build-tarkistukset. Julkaisuputki käynnistyy versiontagista. Kun v1.0:n Windows-TUI-portti on tarkistettu, julkaisun tekijä voi luoda tagin ja lähettää sen:

```powershell
git tag -a v1.0.0 -m "Tasx v1.0.0"
git push origin v1.0.0
```

GitHub Actions rakentaa kuusi käyttöjärjestelmä-/arkkitehtuuripakettia, kerää suoraan linkitettyjen Go-riippuvuuksien lisenssitiedostot `third_party_licenses/`-hakemistoon, tuottaa tarkistussummat ja luo **draft-julkaisun**. Julkaisu pysyy luonnoksena, kunnes pakettien sisältö, lisenssi ja tarkistussummat on hyväksytty GitHubissa.

Tagin push on ulospäin julkaiseva toimi. Älä suorita sitä ilman erillistä käyttäjän pyyntöä.

## Rajat

- TUI:n Windows-käyttö edellyttää oikeaa interaktiivista terminaalia; sitä ei voi vahvistaa buildillä.
- Asennusskripti asentaa vain `tx.exe`-TUI-aliaksen. `tasx`-nimisen vanhan PowerShell-komennon käyttöönotto tai korvaaminen on erillinen päätös.
- Julkaisupaketti ei sisällä käyttäjän tehtävälistoja, asetuksia tai runtime-tietoja.
