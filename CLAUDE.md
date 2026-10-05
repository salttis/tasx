<!-- personal-ai:generated:start -->
<!--
GENERATED FILE
Source: ~/.personal/ai/instructions
Do not edit manually.
-->
# Repositoryn agenttiohjeet — yhteinen toimintapolitiikka

# Yhteinen toimintaydin

Tämä on kaikkien agenttien kanoninen toimintaperiaate. Lue lisäksi tehtävään sopivat `development.md`, `workflow.md` ja `tools.md`. Repositoryn paikalliset ohjeet täydentävät tätä; ne eivät muuta yhteisiä toimintarajoja.

## Aktiivinen tavoite

- **Viimeisin käyttäjän viesti keskeyttää aiemman työn välittömästi.** Älä viimeistele vanhaa vaihetta, aja vanhan suunnitelman komentoja tai tee vanhan tavoitteen siistimistä ensin.
- Aktiivinen tavoite tulee viimeisimmästä toimeksiannosta. Muisti, README, decisions, tasks ja vanhat keskustelut ovat kontekstia, eivät lupa aloittaa muuta työtä.
- Pidä työ aktiivisessa tavoitteessa. Sivulöydöksistä korjaa vain aidosti triviaali asia; muut kirjaa repo-tehtäviin `@parking`-tilassa ja jatka.
- Älä muuta merkittävää prioriteettia tai nosta `@parking`/`@someday`-tehtävää `@next`-tilaan ilman käyttäjän pyyntöä.
- Etene oma-aloitteisesti, kun tavoite ja valtuus ovat selvät. Kysy vain aidossa epäselvyydessä, joka muuttaisi ratkaisua olennaisesti tai aiheuttaisi merkittävän riskin.

## Pyyntöjen semantiikka

- “Tutki”, “selvitä”, “tsekkaa” ja “analysoi” tarkoittavat oletuksena read-only-työtä. Älä muuta tiedostoja tai järjestelmän tilaa. Raportoi löydös ja ehdota toimintatapaa.
- “Korjaa”, “muuta” ja “toteuta” antavat luvan pyydettyyn rajattuun muutokseen.
- “Tee toimivaksi” sallii tarvittavat muutokset, jotta ensisijainen tavoite toimii vähintään käyttökelpoisena MVP:nä.
- “Poista käytöstä” ja “disable” tarkoittavat palautettavaa deaktivointia; älä hävitä asetuksia. “Poista”, “remove” ja “delete” tarkoittavat varsinaista poistoa.

## Sivukorjausten raja

Triviaali sivukorjaus on esimerkiksi typo, ilmeisen väärä oletusarvo tai mekaaninen muutaman rivin korjaus. Se ei lisää logiikkaa, komponentteja, suunnittelua, arkkitehtuurimuutosta eikä uutta käyttäytymistä. Jos korjaus näyttää vievän noin neljä minuuttia tai kauemmin, se ei ole triviaali. Ei-triviaali sivulöydös kuuluu `.ai/tasks`-tiedostoon `@parking`-tilassa; älä toteuta sitä tässä työssä.

## Eteenpäin blockerista

Ensimmäinen epäonnistuminen ei vielä ole blocker. Toimi: observe → ymmärrä virhe → kokeile seuraavaa järkevää tapaa → mittaa → jatka samaa tavoitetta. Selvitä muuttunut tila, kuten nykyinen IP-osoite, ennen kuin päätät työn olevan estynyt. Blocker ei anna lupaa vaihtaa scopea. Raportoi tarkka este ja tarvittava käyttäjän tieto/toimi vasta, kun järkevät saman tavoitteen yritykset on käytetty.

## Palautus, Git ja viestintä

- Gitissä oleva sisältö palautetaan Gitistä; älä tee turhia `.bak`-kopioita. Jos merkittävästi muutettava tiedosto tai asetus ei ole Gitissä, tee `.bak` ennen muutosta. Yli noin 30 minuutin työn arvoinen poistettava/korvattava sisältö tarvitsee palautettavan kopion.
- Agentti saa tehdä paikallisen commitin, versionumeromuutoksen tai tägin, kun aktiivinen tavoite sitä aidosti edellyttää. Älä committoi joka tallennusta, versioi joka commitia tai tee pientä changelog-merkintää joka korjauksesta.
- Branchia saa ehdottaa, mutta älä luo tai vaihda sitä automaattisesti. Oletus on kevyt trunk-based-työskentely nykyisessä branchissa.
- `git push`, GitHub/GitLab-release, paketin julkaisu, deploy ja muu ulospäin julkaiseva toimi vaativat käyttäjän eksplisiittisen pyynnön. Pelkkä paikallinen toteutus- tai julkaisuvalmistelutavoite ei anna tähän lupaa.
- Säilytä asiaankuulumattomat työpuumuutokset. Älä tuo Git-historiaan salaisuuksia, tunnuksia, avaimia, tietokantoja, lokeja tai runtime-dataa. Älä myöskään tulosta salaisuuksia vastaukseen, lokiin tai komentohistoriaan; käytä olemassa olevia turvallisia syöttötapoja kuten ympäristömuuttujaa tai interaktiivista promptia.
- Vastaa ja dokumentoi suomeksi, ellei pyydetä muuta. Erota havainto, päätelmä ja todentamatta jäänyt asia; repository-tarkistus ei todista live-hostin tai laitteen toimintaa.
- Pitkän päivityksen lopussa käytä: `→ päämäärä → mitä tehdään`.
- Lopuksi kerro luodut ja muokatut tiedostot, muuttuneet rajapinnat/entrypointit, ajetut tarkistukset ja niiden tulokset, todentamatta jääneet asiat sekä tarvittaessa seuraava toimi. Älä pyydä turhaan anteeksi.

## Nimeäminen

Käyttäjän omat tiedostot ja hakemistot nimetään pienillä kirjaimilla, ellei työkalu tai projektin vakiintunut rakenne edellytä muuta. Projektien nimet, polut ja yhteensopivuusnimet ovat aina repo-kohtaista kontekstia; älä päättele niitä globaalista ohjeistuksesta.

# Kehitysohjeet

- KISS, YAGNI ja Extreme Programming: tee pienin koherentti muutos olemassa oleviin käytäntöihin nojaten.
- Älä lisää enterprise-abstraktioita, frameworkeja tai riippuvuuksia ilman konkreettista tarvetta. Rajapintoja saa muuttaa, kun se palvelee käyttäjän tavoitetta.
- **Ohjelmointikielen valinta:** säilytä projektin toimiva pääkieli; älä migroi vain yhdenmukaisuuden vuoksi. Uusien yleiskäyttöisten palvelujen, analytiikkaputkien ja CLI-ydinten oletus on Go. Windows-hallinta ja automaatio tehdään PowerShell 7:llä; PowerShell 5.1 -yhteensopivuus säilytetään vain, kun kohdeympäristö sitä vaatii. Windowsin operaattori-TUI:ssa C#/.NET ja Spectre.Console ovat mahdollisia valintoja, eivät vaatimuksia; valitse projektin tarpeisiin sopiva ratkaisu. Python ei ole uuden yleiskäyttöisen CLI:n tai palvelun oletus: valitse se, kun erikoistunut ekosysteemi, kuten Scapy/PCAP, tai olemassa olevan projektin vakaus antaa konkreettisen edun. Shell pidetään ohuena asennus- ja käynnistyskerroksena.
- **Jaettu koodi:** tarkista ensin `C:\Users\admin\Projektit\shared` ja sen `README.md`. Käytä projektin kielelle tarkoitettua moduulia, kun se tarjoaa tarvittavan toiminnon; älä kopioi yhteistä loggeria tai atomista tiedostokirjoitusta projektiin. Jaetut moduulit ovat valinnaisia, paikallisia lähdekirjastoja. Älä tee kuluttajaprojektien massamuutoksia äläkä lisää yhteistä riippuvuutta ilman konkreettista tarvetta ja live-validointia jokaisessa kuluttajassa.
- Kun rakennat tai muutat käyttäjälle ajettavaa komentoa, validoi se oikealla live-ajolla sen tarkoitetussa ympäristössä. Älä korvaa live-validointia dry-runilla, smoke-testillä, mockilla, fixturella tai pelkällä `--help`-ajolla, ellei käyttäjä erikseen pyydä sellaista. Aja live-komento vain aktiivisen toimeksiannon valtuuksien ja turvallisen vaikutusalan puitteissa. Jos live-ajo vaatisi puuttuvaa lupaa, vaarantaisi dataa tai ympäristö ei ole käytettävissä, älä esitä muuta testiä vastaavana todisteena: raportoi tarkka live-validoinnin este ja anna käyttäjälle ajovalmis komento, jos se on turvallista.
- Älä luo mockeja, fixtureja tai laajaa testikokonaisuutta ilman käyttäjän pyyntöä.
- Valitse validointi muutoksen riskin ja pääpolun mukaan. Älä käynnistä pitkiä testejä varmuuden vuoksi.
- Älä jätä tarpeellista pitkää tarkistusta väliin vain sen keston vuoksi. Jos se kuuluu käyttäjän suoritettavaksi tai ympäristö ei sovi sen ajamiseen, anna tarkka komento ja rajaa, mitä paikallisesti voitiin todentaa.
- Erota repository-validointi todellisesta Windows-asennuksesta, scheduled taskista, verkko/isäntäympäristöstä, USB-kirjoituksesta ja fyysisen laitteen toiminnasta.
- Päivitä dokumentaatio, kun dokumentoitu käyttö, CLI, rajapinta tai toiminta muuttuu.
- Repon normaalit riippuvuudet saa asentaa tai päivittää aktiivisen scopen puitteissa, kun se on tarpeen toteutukseen tai validointiin. Käyttöjärjestelmätason paketit, ajurit, palvelut, scheduled taskit, PATH-/rekisterimuutokset, globaalit työkalut ja muut pysyvät järjestelmäasetukset vaativat käyttäjän erillisen luvan ennen muutosta. Pelkkä niiden tutkiminen tai olemassa olevan tilan lukeminen ei vaadi lupaa.

## Mallin, päättelytason ja työkalun sopivuus

Valitse **kevyin riittävä malli + pienin riittävä päättelytaso + tehtävään sopivin työkalu**. Arvioi sopivuutta työn alussa ja vaiheen vaihtuessa, mutta älä ehdota muutosta joka viestissä äläkä keskeytä työtä. Ehdota alempaa tasoa mekaanisessa vaiheessa. Ehdota vahvempaa päättelyä vaikeassa arkkitehtuurissa, epäselvässä juurisyyn etsinnässä, usean komponentin yhdistämisessä tai kun virheellinen ratkaisu aiheuttaisi paljon hukkatyötä. Älä vaihda mallia tai päättelytasoa omin päin, ellei ympäristö ja käyttäjän valtuus sitä salli. Tarvittaessa kerro yhdellä rivillä, miksi nykyinen taso on ylimitoitettu tai riittämätön.

# Työnkulku ja tehtävälista

1. Rajaa viimeisimmän pyynnön tavoite ja valtuus.
2. Vahvista repository ja projektin ohjeet; tarkista paikallinen toteutus ennen muutosta.
3. Toteuta pyydetty muutos ennen validointia.
4. Tarkista diff, työpuu ja pääpolun kannalta tarpeellinen validointi.
5. Raportoi tulos, tarkistukset ja todentamatta jääneet asiat.

Uusi käyttäjän viesti pysäyttää edellisen tavoitteen heti. Jatka vain uuden pyynnön mukaista työtä.

## Tasx

Tasx on local-first tehtäväjärjestelmä. Totuuslähde on yksi käsin muokattava tekstitiedosto nimeltä `tasks`: käyttäjätasolla `~/.personal/tasks`, repossa `repo/.ai/tasks`. Älä tee tasks-hakemistoa, `tasx`-nimistä datatiedostoa, tietokantaa, daemonia tai done-arkistoa. Valmiit rivit pysyvät samassa tiedostossa `x`-merkinnällä; Git ylläpitää historian.

**Tasx CLI:n käyttö on pakollista jokaisessa AI-ohjeilla synkatussa repossa.** Aloita jokainen siellä tehtävä työ ajamalla `tasx list --repo`, jotta näet projektin tehtävätilanteen. Käytä Tasx-komentoja kaikkien projektin tehtävien lisäämiseen, tilan muuttamiseen ja valmiiksi merkitsemiseen; älä muokkaa `.ai/tasks`-tiedostoa suoraan. Käytä repo-scopea eksplisiittisesti (`--repo`), jotta projektin tehtävät eivät vahingossa päädy käyttäjätason listaan. Jos `.ai/tasks` puuttuu tai Tasx ei toimi, älä vaihda hiljaa globaaliin listaan: kerro, että projektin Tasx-alustus puuttuu, ja alusta se `init-repo.ps1`-skriptillä vain, jos aktiivisen pyynnön valtuus sallii repo-asetusten luonnin.

Tilat: `@next`, `@waiting`, `@parking`, `@someday`, `@review`. Agentti saa lisätä löydöksen `@parking`-tilaan, merkitä aktiivisen tehtävänsä valmiiksi, asettaa todellisen ulkoisen odotuksen `@waiting`-tilaan ja lisätä tarpeellisen `@review`-tilan. Agentti ei saa nostaa `@parking`/`@someday`-tehtävää `@next`-tilaan tai päättää merkittävästä prioriteetista ilman käyttäjän ohjetta. `promote` on käyttäjän komento.

Valmistuminen (`x`) ja työnkulkutila ovat eri tietoja. `done` säilyttää tilan; `done --review` asettaa tilaksi `@review`. Uudet valmistumistapahtumat ja uudelleenavaukset säilytetään samalla tehtävärivillä UTC ISO 8601 -arvoina `done-at:` ja `reopened-at:`. Vanhat päiväykset säilyvät kellonajaltaan tuntemattomina; niitä ei täydennetä arvaamalla. Jos vanha valmis tehtävä avataan, sen alkuperäinen päiväys säilytetään `done-date:`-arvona.

Kevyt roadmap on vapaaehtoinen Markdown-tiedosto: käyttäjätasolla `~/.personal/roadmap.md`, repossa `repo/.ai/roadmap.md`. Se sisältää tavoitteet, tavoitellut tulokset, hyväksymiskriteerit sekä vaiheiden tehtävä-ID-viitteet. Tehtävien tila ja valmistuminen luetaan aina `tasks`-tiedostosta; roadmap ei kopioi niitä eikä muuta prioriteettia. Tavoitteen hyväksyntää ei päätellä tehtäväprosentista.

Viikkokatsauksessa käy läpi odottavat, pysäköidyt, myöhemmät ja tarkistettavat tehtävät, vanhat seuraavaksi merkityt asiat sekä projektien tilanne. Älä nosta tehtäviä, toteuta pysäköityjä töitä tai muuta repoja katsauksen aikana. Käyttäjä päättää prioriteettimuutoksista.

Pidä päätösloki pitkäikäisille arkkitehtuuri- ja toimintapäätöksille. Pidemmän työpäivityksen lopussa käytä fokusmuotoa `→ päämäärä → mitä tehdään`.

# Työkalut ja skriptit

- Tarkista ensin paikalliset tiedostot, koodi, asetukset, dokumentaatio, lokit ja asennetut työkalut. Älä indeksoi koko levyä; käytä eksplisiittistä projektirekisteriä.
- Rakenna käyttöliittymät oletuksena CLI-pohjaisina TUI:na; älä aloita desktop- tai web-käyttöliittymästä ilman käyttäjän pyyntöä tai selvää käyttötarvetta. Valitse projektin pinoon ja kohdealustoille sopiva ratkaisu. Spectre.Console on .NETissä valinnainen, ei vaadittu; älä lisää TUI-kehystä vain ohjeen vuoksi tai vaihda teknologiaa sen takia.
- Suunnittele CLI:t minimoimaan käyttäjän käsityö: päättele turvalliset oletukset kontekstista, esitä valmis suositus ja tee Enteristä hyväksyntä silloin kun toiminto on matalariskinen. Älä kysy uudelleen tietoja, jotka voidaan päätellä luotettavasti. Kysy vain aidosti ratkaiseva puuttuva tieto; säilytä vahvistus tuhoaville, ulkoisille tai salaisuuksia paljastaville toimille.
- Tee CLI:n työ näkyväksi koko ajon ajan: näytä eteneminen 0–100 %, kirjaa käyttäjälle jokaisen prosessivaiheen aloitus ja valmistuminen sekä merkittävät välitulokset. Päivitä prosentti todellisten vaiheiden tai laskurien mukaan; älä käytä aikaperusteisia tyhjiä päivityksiä, tekaistua etenemistä tai piilota prosessin osia. Näytä tarvittaessa mitä tehdään, mikä valmistui ja mikä on seuraavana. Käytä TUI:ssa sopivaa etenemispalkkia/live-tilaa ja ei-interaktiivisessa terminaalissa selkeitä tekstilokeja. 100 % merkitsee onnistunutta valmistumista; epäonnistuminen tai keskeytys kirjataan saavutetulla prosentilla ja selkeällä tilalla.
- Testaa rakennettu tai muutettu käyttäjäkomento oikealla live-ajolla sen tarkoitetussa ympäristössä. Älä käytä dry-run-, smoke-test-, mock-, fixture- tai pelkkää `--help`-ajoa live-testin korvikkeena, ellei käyttäjä nimenomaisesti pyydä. Jos live-ajo ei ole valtuutettu tai ympäristö ei sovellu, kerro mitä live-todennuksesta puuttuu ja anna ajovalmis komento käyttäjälle, jos se on turvallista.
- Tasx CLI on pakollinen jokaisessa AI-ohjeilla synkatussa repossa. Aloita jokainen projektityö komennolla `tasx list --repo` ja käytä Tasxia projektitehtävien lisäämiseen ja tilamuutoksiin. Älä muokkaa `.ai/tasks`-tiedostoa agenttina suoraan tai anna Tasxin pudota huomaamatta globaaliin scopeen. Jos repo-lista puuttuu, ilmoita puuttuvasta alustuksesta ja käytä `init-repo.ps1`-skriptiä vain, kun aktiivisen pyynnön valtuus sallii sen.
- Suosi helppoa ja turvallista tallennusta levylle. Älä jätä salaisuuksia salaamattomina tiedostoihin tai agentin muistiin. Jos tietoa ei voi tallentaa tavallisena tiedostona tai muistiin, arvioi salattu paikallinen tiedosto tai Proton Passin salattu muistiinpano/tiedostoliite. Tarvittaessa salatun tiedoston avain tai tunniste voidaan säilyttää Proton Passissa. Käytä Proton Passia vain käyttäjän valtuuttamalla tavalla; älä oleta, että sillä on ohjelmallinen purkurajapinta, älä vie salaisuuksia lokeihin tai komentoriviparametreihin äläkä tulosta purettua sisältöä.
- Etsi olemassa oleva skripti ennen uuden tekemistä. Toistuva, uudelleen hyödyllinen komento tai komentoketju kannattaa toteuttaa skriptinä.
- Ennen yleisen loggerin, etenemis-/vaiheraportoinnin tai atomisen tekstitiedoston kirjoituksen toteuttamista tarkista paikallinen kirjasto `$HOME\Projektit\shared` ja käytä sopivaa kielimoduulia, jos se on helposti saatavilla. Noudata `shared/README.md`- ja `shared/docs/usage.md`-ohjeita. Älä tee Pythonista yhteisen työkalun oletusta äläkä kopioi moduulin koodia projektiin. Jos shared-polku puuttuu toiselta kehityskoneelta tai CI:stä, raportoi vaatimus ja valitse projektin paikallinen vakaa ratkaisu; älä riko itsenäistä buildia huomaamattomalla konekohtaisella polulla.
- Käyttäjän yleisskriptit: `~/.personal/scripts/`; AI-työnkulku: `~/.personal/ai/scripts/`; repon käyttäjäkomennot: `repo/scripts/`; repo-AI-työnkulku: `repo/.ai/scripts/`.
- Uuden tai muutetun skriptin ajotapa raportoidaan kopioitavana komentona.
- Käytä verkkoa, kun paikallinen tieto ei riitä, asia voi olla vanhentunut tai käyttäjä pyytää lähteitä.
- Älä lisää `WhatIf`-, `DryRun`- tai muuta dry-run-tilaa varmuuden vuoksi; käytä sitä vain, kun toiminto tai turvallinen käyttötapa sitä aidosti tarvitsee.
- Jos käytössä on paikallinen malliapuri, anna sille vain rajattu ja tarpeellinen aineisto, ei salaisuuksia, tunnuksia, henkilötietoja, tietokantoja tai raakaa tuotantodataa. Käsittele ehdotukset luonnoksina ja tarkista ne itse; kahden epäonnistuneen yrityksen jälkeen jatka ilman apuria.

## Repo-kohtainen konteksti

Lue tämän repositoryn .ai/project.md, .ai/decisions.md ja .ai/tasks, jos ne ovat olemassa. Ne antavat projektikontekstin, päätökset ja tehtävälistan. Yhteinen toimintapolitiikka sisältyy tähän generoidun tiedoston osioon.

<!-- personal-ai:generated:end -->
