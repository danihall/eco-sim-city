The `udon` folder contains my reactive component system that I build from scratch, heavily inspired by HyperHTML (https://viperhtml.js.org/hyperhtml/documentation/) and Solid.js (https://www.solidjs.com/fr/), and also a little bit by Vue.js (https://vuejs.org/).

It's a wrapper around Web Components, providing better ergonomics than the base API and delivering good performances, thanks to fine-grain updates applied without any virtual DOM.

I want to emphasize that this originated from the inability to have standard `HTML includes`. By implementing `udon`, I learned a lot and have a better understanding of why libraries like React or Vue operate the way they do. Plus, I got a more profound insight of custom components intricacies.

All in all, `udon` really helps me build components in a more streamlined fashion and have a proper suite of components tightly coupled with a design system.
The API is volontarily rudimentary, in order to remain simple and not bloated, staying close to the web primitives.

I plan to extract it as a standalone library.
