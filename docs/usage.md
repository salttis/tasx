# Käyttö

## Rakentaminen ja käynnistys

Repojuuresta:

```powershell
go run ./cmd/tasx --help
go run ./cmd/tasx list
go run ./cmd/tasx list tasx
go run ./cmd/tasx roadmap
go run ./cmd/tasx roadmap tasx
go run ./cmd/tasx add "Toteuta uusi toiminto #development @next"
go run ./cmd/tasx add "Toteuta projektin toiminto #ui @next" --projekti tasx
go run ./cmd/tasx done 2
go run ./cmd/tasx done tasx-2
go run ./cmd/tasx reopen 2
go run ./cmd/tasx status 2 waiting
go run ./cmd/tasx archive 2
go run ./cmd/tasx config
go run ./cmd/tasx tui
```

Go-binäärin `tasx`-nimi listaa ilman argumentteja. `tx`-nimellä käynnistetty sama binääri avaa TUI:n ilman argumentteja. Molemmat binäärinimet hyväksyvät eksplisiittiset alikomennot, kuten `list` ja `tui`.

Asenna `tx` nykyisen käyttäjän komentohakemistoon:

```powershell
pwsh -NoProfile -File .\scripts\install.ps1
```

Asennusskripti rakentaa `dist\tx.exe`-binäärin ja kopioi sen `~\.personal\scripts`-hakemistoon. Se ei muuta PATH-asetusta. Julkaisuversion voi antaa `-Version`-valitsimella; jos aiempi `tx.exe` on olemassa, se säilytetään palautuskopiona. Lisätiedot: [julkaisuopas](./release.md).

## Nykyiset Go-komennot

| Komento | Toiminta |
|---|---|
| `tasx list [projekti] [--state all\|open\|done]` | Ilman projektia listaa ajokontekstin tehtävälistan; projektin nimellä listaa rekisteröidyn projektin mistä tahansa. Tuloste on ihmisluettava teksti. |
| `tasx roadmap [projekti] [--repo tai --global]` | Näyttää valitun scopen roadmapin tavoitteet, hyväksymiskriteerit, vaiheet ja tehtäväriippuvuudet tekstinä. |
| `tasx add <kuvaus> [--projekti <nimi>] [--tag <tagi>] [--status <tila>]` | Lisää nykyiseen projektiin tai `--projekti`-valitsimella rekisteröityyn projektiin. Antaa seuraavan juoksevan `$numero`-ID:n. Tagit ja tilan voi kirjoittaa myös kuvaukseen, kuten `#ui @next`. |
| `tasx done <numero\|projekti-numero>` / `tasx reopen <numero\|projekti-numero>` | Muuttaa checkbox-valmiutta ja päivittää saman rivin aikaleiman. Numero yksin valitsee ajokontekstin projektin; projektinimi-numero toimii mistä tahansa. Tehtävä jää tiedostoon. |
| `tasx status <numero\|projekti-numero> <tila\|none>` | Asettaa tai poistaa työnkulkutilan (`next`, `waiting`, `parking`, `someday`, `review`). |
| `tasx archive <numero\|projekti-numero>` | Käyttäjän erikseen käynnistämä siirto `Arkisto:`-osioon; siirtää valitun tehtävän alitehtävineen ja säilyttää ID:t. |
| `tasx tui [--repo tai --global]` | Avaa TUI:n tehtävien tarkasteluun ja muokkaukseen. |
| `tasx config` | Tulostaa käytössä olevat asetukset ja asetustiedoston polun tekstinä. |
| `tasx version` | Tulostaa version (kehitysbuildissa `dev`). |

Oletuksena Tasx päättelee scopen ajokontekstista: repossa käytetään sen `.ai/tasks`-listaa, muualla käyttäjälistaa. `--repo` ja `--global` ovat eksplisiittisiä ylikirjoituksia; ne eivät ole yhdistettävissä. Projektin nimi `list`-komennossa, ID:n projektietuliitteessä tai `add --projekti`-valitsimessa hakee scopen projektirekisteristä.

Tasx tukee vain Markdown-checkbox-formaattia. Vanha formaatti ja migraatiokomennot eivät ole tuettuja. CLI ei tuota JSON-tulosteita.

## TUI-näppäimet

| Näppäin | Toiminta |
|---|---|
| `Tab`, `h`, `l` | Vaihda paneelia. |
| `j`, `k`, nuolet | Selaa scopen tai tehtävien listaa. |
| `/` | Aloita haku ID:stä tai tehtävän kuvauksesta. |
| `f` | Vaihda suodatinta: kaikki, avoimet, valmiit. |
| `a` | Lisää tehtävä; Enter tallentaa ja Esc peruuttaa. |
| `Space`, `d` | Vaihda valitun tehtävän valmius. |
| `s` | Kierrätä työnkulkutilaa: ei tilaa, `next`, `waiting`, `parking`, `someday`, `review`. |
| `r` | Lue valitun scopen tehtävälista ja roadmap uudelleen. |
| `?` | Näytä näppäinohje. |
| `q`, `Ctrl+C` | Sulje TUI. |

TUI:n tehtävämuutokset käyttävät samaa tallennuskerrosta kuin CLI. Alle 90 sarakkeen terminaalissa kerrallaan näytetään valittu paneeli. Käyttöä aidossa Windows-terminaalissa ei ole vielä vahvistettu.
