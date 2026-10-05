# Tasx-projektikonteksti

- Tavoite: local-first tehtävienhallinta Go CLI:llä ja Lazygit-tyylisellä TUI:lla.
- Teknologiat: Go 1.27+, Cobra CLI:lle, Bubble Tea ja Lipgloss TUI:lle, TOML käyttäjäasetuksille.
- Käynnistys: `go run ./cmd/tasx`; `go run ./cmd/tasx tui`; `tx`-niminen build käynnistää TUI:n oletuksena.
- Data: käyttäjätason `~/.personal/tasks`; repo-scope `.ai/tasks`; ei tietokantaa eikä daemonia. `--repo` vaatii repo-listan, kun taas `default_scope = "auto"` voi valita käyttäjälistan, jos repo-listaa ei ole.
- Käyttäjäasetukset: `~/.tasxrc`, polun voi vaihtaa `TASX_CONFIG`-muuttujalla. `TASX_HOME` valitsee käyttäjädatan juuren.
- Tärkeät rajat: Go tukee vain Markdown-checkbox-formaattia. Vanhaa formaattia tai muunnoksia ei tueta. CLI:n tuloste on tekstimuotoista, ei JSON.
- Jaetut kirjastot: `C:\Users\admin\Projektit\shared\go` on valinnainen lähdemoduuli; Tasxin itsenäisen buildin ei pidä riippua konekohtaisesta `replace`-polusta.
- Nykytila (2026-10-05): Go CLI sisältää tehtävien lisäys-, valmistumis-, uudelleenavaus-, tila- ja käyttäjän käynnistämän arkistointikomennon sekä `roadmap`-tekstikomennon. ID merkitään tehtäväriville `$numero`-muodossa ja kvalifioidaan komennoissa projektilla (`projekti-numero`); ajokontekstin projektissa riittää numero. Oletusscope havaitaan repo-/käyttäjäkontekstista; rekisteröidyn projektin voi valita `list`-komennolla tai `add --projekti`-valitsimella. Arkistointi säilyttää tehtävän ja alitehtävien ID:t eikä käynnisty automaattisesti. TUI:n kolmen paneelin näkymässä voi lisätä tehtäviä, vaihtaa valmiutta ja kierrättää työnkulkutilaa. Roadmap näytetään yksityiskohtapaneelissa. Interaktiivista käyttöä ei ole vielä vahvistettu aidossa Windows-terminaalissa.
- Suunnitelma ja dokumentaatio: `docs/README.md`, `docs/project-plan.md`, `docs/roadmap.md`; koneellisesti luettava roadmap on `.ai/roadmap.md`.
- Projektin tavoite- ja vaihekuvaus pidetään `.ai/roadmap.md`-tiedostossa; Tasxin tehtävälista poistetaan vanhana käyttäjän pyynnöstä.
- Validointi: Go build ja CLI:n live-ajot muutetuille toiminnoille; interaktiivinen TUI tarkistetaan oikeassa terminaalissa.
