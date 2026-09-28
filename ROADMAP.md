# 🚀 AstroLeap: Feuille de Route & Évolutions Futures (Backlog)

Ce document conserve les idées validées par l'équipe pour un développement ultérieur.

---

## 🤖 1. Évolution Active du Drone Compagnon
*Actuellement cosmétique et suiveur (avec animations idle/float/hurt intégrées).*

* **Balayage de Proximité (*Secret Proximity Sonar*)** :
  * Le drone émet un bip sonore aigu et des étincelles cyan lorsqu'un Cristal d'énergie caché, une Médaille d'or ou un Terminal Holographique se trouve dans un rayon de 60 px.
  * Déclenche une pulsation visuelle autour du joueur pour orienter la recherche dans les zones secrètes.
* **Bouclier d'Urgence (*Pulse Shield*)** :
  * Capacité activable (ex: double-tap Bas ou bouton dédié) consommant 50% de la jauge de carburant jetpack.
  * Génère un champ de force sphérique bleu translucide de 1,5 seconde qui absorbe un coup mortel ou repousse les projectiles ennemis mineurs.
* **Hacking de Barrières de Sécurité** :
  * Dans certaines salles verrouillées, l'astronaute peut positionner son drone devant une console pour désactiver temporairement un champ de force laser.

---

## 🔊 2. Ambiance Sonore & Synthèse Procédurale (SFX / Radio)
*Actuellement géré par le moteur audio sans dépendance externe (`internal/audio`).*

* **Effets Sonores Radio / Codex** :
  * Synthèse en temps réel d'un jingle d'ouverture de transmission (accord synthétique descendant + bip de modulation de fréquence FM).
  * Effet sonore de déchiffrement pas-à-pas (effet télex / frappe numérique rapide) lors de la collecte d'un terminal en jeu.
* **Bruitages de Gravité Zéro & Sas** :
  * Dépressurisation de sas (bruit blanc filtré passe-bas avec enveloppe exponentielle).
  * Son d'alerte de verrouillage de cible de la grille A.E.G.I.S.
