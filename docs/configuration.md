# Asetukset

Tasx lukee valinnaisen käyttäjäkohtaisen TOML-tiedoston polusta `~\.tasxrc`. Puuttuva tiedosto käyttää sisäänrakennettuja oletuksia. Virheellinen TOML, virheellinen arvo tai tuntematon asetus ilmoitetaan virheenä.

```toml
personal_dir = "" # tyhjä arvo käyttää ~/.personal-hakemistoa
default_scope = "auto" # auto, repo, global
state = "all" # all, open, done
```

## Asetukset

| Avain | Oletus | Merkitys |
|---|---|---|
| `personal_dir` | `~/.personal` | Käyttäjän Tasx-tietojen ja projektirekisterin juurihakemisto. |
| `default_scope` | `auto` | `auto` käyttää nykyistä repo-scopea, kun `.ai/tasks` löytyy; muussa tapauksessa käyttäjäscopea. `repo` vaatii repo-scopen. `global` valitsee käyttäjäscopen. |
| `state` | `all` | Listauksen ja TUI:n oletussuodatin: kaikki, avoimet tai valmiit. |

## Ympäristömuuttujat

- `TASX_CONFIG` vaihtaa asetustiedoston sijainnin.
- `TASX_HOME` asettaa käyttäjädatan juuren ja ohittaa `personal_dir`-arvon.

## Valintojen etusija

1. Komennon `--repo` tai `--global` valitsee scopen eksplisiittisesti ja ohittaa `default_scope`-oletuksen.
2. CLI-valitsin `--state` ohittaa asetustiedoston oletuksen.
3. Muussa tapauksessa käytetään `~/.tasxrc`-arvoja tai sisäänrakennettuja oletuksia.

`--repo`-valinta epäonnistuu, jos nykyinen hakemisto ei ole Git-repossa tai repo-lista puuttuu. `default_scope = "auto"` saa siirtymävaiheen Go-versiossa käyttää käyttäjälistaa, kun repo-listaa ei ole; tätä ei pidä sekoittaa eksplisiittiseen `--repo`-valintaan. Sama lippuvalinta koskee CLI:tä ja TUI:n käynnistystä.

Tulostetut asetukset:

```powershell
tasx config
```

Komento näyttää asetukset tavallisena tekstinä.
