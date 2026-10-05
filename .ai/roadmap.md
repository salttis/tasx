# Tasx-roadmap

## goal:go-migration | Luotettava Go CLI ja TUI
Outcome: Tasx toimii itsenäisenä local-first Go-ohjelmana, tukee Markdown-tehtäväformaattia, CLI:tä ja kirjoittavaa TUI:ta.
Acceptance criteria:
- [x] Markdown-checkbox-formaatti lukee osiot, hierarkian, kommentit ja metatiedot.
- [x] Vanhaa tehtäväformaattia tai migraatiota ei tueta.
- [x] Kirjoittavat operaatiot säilyttävät muokkaamattomat rivit ja metatiedot sekä suojaavat tiedoston virheiltä.
- [x] Roadmap-rajapinta näyttää tavoitteet, hyväksymiskriteerit, vaiheet ja tehtäväriippuvuudet tekstinä.
- [x] CLI ei tarjoa JSON-ulostuloa.
- [x] TUI tarjoaa tehtävien käsittelyn samalla tallennuskerroksella kuin CLI.
- [ ] TUI:n interaktiivinen käyttö vahvistetaan aidossa Windows-terminaalissa ennen v1.0-julkaisun julkaisemista.
- [ ] `tx`- ja `tasx`-käynnistysmigraatio tehdään vasta yhteensopivuuden tarkistamisen jälkeen.
- [x] v1.0-julkaisun asennus tukee versionumeroa, säilyttää edellisen binäärin ja tarjoaa palautuksen.
- [x] v1.0-julkaisupaketit ja tarkistussummat rakentuvat Git-tagista, ja julkaisu luodaan ensin draftina.

### phase:documentation | Projektin dokumentaatio

### phase:markdown-format | Markdown-tehtäväformaatti

### phase:write-cli | Säilyttävät kirjoitusoperaatiot ja CLI

### phase:scopes-roadmap | Scope-koosteet ja roadmap-rajapinta

### phase:tui-workflows | TUI:n tehtävätoiminnot

### phase:command-migration | Komentomigraatio ja vakautus
