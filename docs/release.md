# Julkaisu

## v1.0.0 sisältö

v1.0.0 on paikallinen Go-sovellus, joka tarjoaa Markdown-checkbox-tehtävien CLI-hallinnan ja Bubble Tea + Lipgloss -TUI:n. Julkaisu sisältää Windows-, Linux- ja macOS-binäärit `amd64`- ja `arm64`-arkkitehtuureille, `tasx`- ja `tx`-nimiset käynnistimet, sovellus- ja kolmannen osapuolen lisenssitiedot sekä SHA-256-tarkistussummat. Tasxin oma lisenssi on MIT.

Vanhaa PowerShell-komentopolkua ei korvata automaattisesti.

## Minor- ja major-julkaisu

Aja PowerShell 7:llä repojuuresta:

```powershell
.\scripts\release.ps1 minor
.\scripts\release.ps1 major
```

`minor` kasvattaa uusimman vakaan tagin minor-version ja nollaa patch-version; `major` kasvattaa major-version ja nollaa minor- ja patch-versiot. Esimerkiksi `v1.0.0`-tagista seuraavat versiot ovat `v1.1.0` ja `v2.0.0`.

Julkaisu edellyttää `git`, `go`, `gh` ja `tar` -komentoja, GitHub CLI -kirjautumista `salttis/tasx`-repooon, puhdasta työpuuta sekä ajan tasalla olevaa `main`-branchia. Tee julkaistavat muutokset commitiksi ja puske ne `origin/main`-haaraan ennen julkaisua. Julkaisun työvaiheet:

1. Etsi uusin vakaa semanttinen tagi ja laske uusi versio. Jos tagiin liittyvä julkaisu on jo kesken, sama komento jatkaa sitä.
2. Aja `gofmt`-tarkistus, `go test ./...`, `go vet ./...` ja build kaikille kuudelle Windows-, Linux- ja macOS-kohteelle (`amd64` ja `arm64`).
3. Pyydä kirjoittamaan ehdotettu tagi. Väärä vastaus keskeyttää julkaisun ennen tagin puskemista.
4. Luo annotatoitu tagi ja puske se `origin`-remoteen. Release workflow tekee tarkistukset ja paketit GitHub Actionsissa.
5. Jos GitHub ei osoita workflow'lle runneria viidessä minuutissa, skripti peruu jonoon jääneen ajon ja paketoi kuusi arkistoa paikallisesti erilliseen väliaikaishakemistoon. Paikallinen varatapa rakentaa binäärit, kerää kolmannen osapuolen lisenssit ja luo SHA-256-listan.
6. Lataa draftin assetit, tarkistaa kuuden paketin sisällön, lisenssitiedostot ja SHA-256-tiivisteet, ja julkaisee draftin julkiseksi vasta kaikkien tarkistusten onnistuttua.

Julkinen release on ulkoinen ja pysyvästi käyttäjille näkyvä toiminto. Tarkista syötetty versio ennen vahvistusta. Workflow'n epäonnistuminen tai assettien tarkistusvirhe keskeyttää julkistamisen; jo pushettu tagi ja mahdollinen draft jäävät tarkastettaviksi. Jos skripti keskeytyy muuten kesken julkaisun, aja sama `major`- tai `minor`-komento jatkaaksesi sitä. Selvitä workflow'n epäonnistuminen ennen kuin yrität jatkaa.

## Versioitu paikallisasennus

PowerShell 7 -repojuuresta:

```powershell
pwsh -NoProfile -File .\scripts\install.ps1 -Version 1.0.0
```

Asennus rakentaa `dist\tasx.exe`- ja `dist\tx.exe`-binäärit ja kopioi ne `~\.personal\scripts`-hakemistoon. `tasx.exe` käynnistää CLI:n ja `tx.exe` TUI:n. Asennus ei muuta PATHia. Jos kohteessa oli jo binääri, edellinen tallennetaan nimellä `<komento>.exe.previous`; jos varmuuskopio oli jo olemassa, uusi saa aikaleimatun nimen. Jos `tasx.exe` puuttui ennen asennusta, asennin kirjaa tämän, jotta palautus poistaa uuden binäärin.

```powershell
pwsh -NoProfile -File .\scripts\install.ps1 -Rollback
```

Palautus käsittelee `tasx.exe`- ja `tx.exe`-binäärit erikseen ja palauttaa kummankin edellisen version, jos se on olemassa.

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

## GitHub-julkaisun toteutus

CI suorittaa pushissa, pull requestissa ja manuaalisesti muotoilu-, `go vet`- ja build-tarkistukset. `release.ps1` luo versiontagin vasta tarkistusten jälkeen:

```powershell
.\scripts\release.ps1 minor
```

GitHub Actions rakentaa kuusi käyttöjärjestelmä-/arkkitehtuuripakettia, kerää suoraan linkitettyjen Go-riippuvuuksien lisenssitiedostot `third_party_licenses/`-hakemistoon, tuottaa tarkistussummat ja luo draft-julkaisun. Julkaisuskripti tarkistaa ladatut assetit ennen draftin muuttamista julkiseksi. Jos Actions-runnerin jono ylittää viisi minuuttia, skripti käyttää paikallista paketointia.

Tagin push ja releasen julkistaminen ovat ulospäin julkaisevia toimia. Käytä release-skriptiä vain, kun uusi versio on hyväksytty julkaistavaksi.

## Rajat

- Interaktiivinen TUI-testaus ei kuulu automaattisen release-skriptin portteihin.
- Asennusskripti asentaa `tasx.exe`-CLI:n ja `tx.exe`-TUI-aliaksen. Se ei muuta eikä korvaa erillistä `tasx.ps1`-komentoa.
- Julkaisupaketti ei sisällä käyttäjän tehtävälistoja, asetuksia tai runtime-tietoja.
