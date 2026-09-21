# 🚀 AstroLeap: Lunar Odyssey

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Engine](https://img.shields.io/badge/Engine-Ebitengine%20v2-db5858)](https://ebitengine.org)
[![Platform](https://img.shields.io/badge/Platform-Desktop%20%7C%20WebAssembly-brightgreen)](https://ebitengine.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**AstroLeap: Lunar Odyssey** est un jeu de plateforme 2D rétro en basse gravité développé en Go avec **Ebitengine v2**. 

Inspiré des classiques rétro comme *Super Mario Bros*, *Moon Patrol* et les mécaniques de vol de *Jetpack Joyride*, vous incarnez un astronaute explorant des environnements lunaires et extraterrestres hostiles pour réparer un vaisseau mère et rentrer sur Terre.

![AstroLeap Logo](logo.png)

---

## 📖 Présentation du jeu

### 🌌 Synopsis & Objectif
À la suite d'une avarie critique, vous devez traverser 5 secteurs planétaires et stations spatiales abandonnées. Maîtrisez la faible gravité lunaire et votre propulseur dorsal (jetpack), éliminez ou évitez les créatures extraterrestres, ramassez les cristaux d'énergie, découvrez des armes secrètes et activez le moteur de saut supraluminique (*Warp Drive*) pour vous échapper !

### 🗺️ Les 5 Secteurs
1. **Secteur 1 : Avant-poste Lunaire (*Lunar Outpost*)**  
   Prise en main des sauts flottants en basse gravité, gestion du jetpack, déblocage du **Laser Blaster** et sauvetage initial.
2. **Secteur 2 : Crête de Phobos (*Phobos Ridge*)**  
   Canyons martiens escarpés, fosses d'acide, plateformes mobiles synchronisées et déblocage de l'arme secrète **Nova Cannon**.
3. **Secteur 3 : Cœur de Glace d'Europe (*Europa Ice Core*)**  
   Cavernes glaciaires d'Europe avec combat de boss contre l'**Overlord Mech**.
4. **Secteur 4 : Couloir du Vaisseau Mère (*Mothership Corridor*)**  
   Intérieur d'un vaisseau mère humain en confinement d'urgence. Récupérez la **carte d'accès de sécurité** (*Keycard*) pour déverrouiller le sas blindé.
5. **Secteur 5 : Salle du Réacteur (*Reactor Bay*)**  
   Cœur du réacteur en surchauffe critique. Récupérez les **3 cœurs de réparation Warp** à travers les pistons et les cheminées de plasma pour initier le saut de retour vers la Terre !

---

## ✨ Fonctionnalités & Mécaniques

- **Physique en basse gravité ($g \approx 0.18$)** : Sauts amples et planés offrant un contrôle aérien précis.
- **Propulseur dorsal (Jetpack)** : Maintenez le bouton de saut en l'air pour planer et franchir de larges précipices. La jauge de carburant se recharge automatiquement au sol ou via des capsules d'énergie.
- **Élimination des ennemis & Rebond** : Écrasez les slimes et yeux volants d'un coup de botte bien placé pour rebondir plus haut.
- **Arsenal d'armes déblocables** :
  - 🔹 **Laser Blaster** : Rayons photoniques cyan traversant les aliens horizontalement.
  - 🔸 **Nova Cannon** : Projectiles à ricochet rebondissant sur les murs et le sol.
  - 🔄 **Changement d'arme à la volée** via une simple touche.
- **Système de records & Speedrun** :
  - Sauvegarde automatique locale de vos meilleurs scores et temps (`~/.astroleap_records.json`).
  - Système d'évaluation de mission (Rangs **S**, **A**, **B**, **C**).
- **Rendu Rétro Pixel Art 16-bit** : Résolution virtuelle native $320 \times 180$ avec mise à l'échelle plein écran nette, effets de particules procéduraux, parallaxe stellaire et filtre CRT rétro.
- **Audio & Musique Intégrés** : Synthèse audio procédurale multi-voix et effets sonores générés en temps réel sans dépendance externe.

---

## 🕹️ Contrôles du jeu

Le jeu prend en charge le **clavier**, les **manettes** (Xbox / PlayStation / génériques) ainsi que les **écrans tactiles** (en version Web).

| Action | Clavier | Manette (Gamepad) | Souris / Tactile (Web) |
| :--- | :--- | :--- | :--- |
| **Déplacement Gauche / Droite** | `A` / `D` ou `←` / `→` | Stick Gauche / Croix directionnelle | Zones Gauche / Centre |
| **Saut** | `Espace` / `W` / `↑` | Bouton Sud (`A` / `✕`) | Bouton Droit |
| **Propulseur (Jetpack)** | Maintenir `Espace` / `W` en l'air | Maintenir Bouton Sud | Maintenir Bouton Droit |
| **Tir** | `J` / `Z` / `F` / Clic Gauche | Bouton Ouest (`X` / `◻`) | Clic / Tap écran |
| **Changer d'arme** | `Q` / `Tab` / `E` | Bouton Nord (`Y` / `△`) / `R1` | Tap icône Arme HUD |
| **Pause** | `Échap` / `P` | Bouton Start | Bouton Pause HUD |
| **Recommencer** | `R` (sur Game Over) | Select / Back | Tap écran |

---

## 🚀 Démarrage Rapide (Getting Started)

### 📋 Prérequis
- **Go 1.22 ou supérieur** ([Télécharger Go](https://go.dev/dl/))
- **Système d'exploitation** : Linux, macOS ou Windows
- *(Linux uniquement)* : Les bibliothèques graphiques et audio standard (`libasound2-dev`, `libgl1-mesa-dev`, `libxcursor-dev`, etc. généralement déjà présentes sur Ubuntu/Debian/Fedora/Arch).

### 📥 1. Cloner le projet

```bash
git clone git@github.com:maximelarrieu/astroleap.git
cd astroleap
```

Téléchargez les dépendances Go :
```bash
go mod download
```

---

### 🎮 2. Lancer le jeu

Vous avez deux façons de jouer : en **application native desktop** ou dans votre **navigateur web**.

#### Option A : Version Desktop (recommandée)

##### Méthode 1 : Lancement direct
```bash
go run ./cmd/game
```

##### Méthode 2 : Compilation d'un binaire exécutable
Pour compiler un binaire autonome :

```bash
# Sur Linux / macOS :
go build -o astroleap ./cmd/game
./astroleap

# Sur Windows (PowerShell / CMD) :
go build -o astroleap.exe ./cmd/game
.\astroleap.exe
```

---

#### Option B : Version Web (WebAssembly)

Le jeu fonctionne également dans n'importe quel navigateur moderne grâce à WebAssembly. Un serveur web HTTP léger est inclus dans le projet.

1. **Lancez le serveur local :**
   ```bash
   go run ./cmd/server
   ```
   *(ou compilez le serveur : `go build -o server ./cmd/server && ./server`)*

2. **Ouvrez votre navigateur :**
   Rendez-vous sur 👉 **[http://localhost:8080](http://localhost:8080)**

*(Optionnel) Si vous modifiez le code du jeu et souhaitez regénérer le binaire WASM :*
```bash
GOOS=js GOARCH=wasm go build -o web/game.wasm ./cmd/game
```

---

## 🧪 Lancer les tests

Pour exécuter l'ensemble des tests unitaires du projet :

```bash
go test -v ./...
```

---

## 📁 Structure du projet

```text
.
├── cmd/
│   ├── game/             # Point d'entrée pour la version Desktop native
│   └── server/           # Serveur HTTP local pour la version WebAssembly
├── internal/
│   ├── assets/           # Générateur de sprites, textures et atlas procéduraux
│   ├── audio/            # Moteur sonore DSP & synthèse de musique de fond
│   ├── entity/           # Entités : Joueur, Ennemis, Cristaux, Boss, Missiles
│   ├── game/             # Boucle principale ebiten.Game
│   ├── input/            # Gestionnaire multi-périphériques (Clavier, Manette, Tactile)
│   ├── level/            # Niveaux, tuiles et décors en parallaxe
│   ├── physics/          # Collisions AABB sweep et gravité lunaire
│   ├── records/          # Persistance des scores et chronos
│   ├── shaders/          # Shaders Kage (effet tube cathodique CRT)
│   ├── state/            # Machine à états (Intro, Titre, En jeu, Victoire, Game Over)
│   └── ui/               # Interface utilisateur, affichage HUD rétro
├── web/                  # Fichiers statiques WebAssembly (HTML, JS runner, WASM)
├── GDD.md                # Document de Game Design complet
└── README.md             # Documentation et guide de démarrage
```

---

## 📜 Licence

Ce projet est distribué sous licence MIT. Amusez-vous bien dans l'espace ! 🚀🌕