## IMPORTANT
Inside workers `.ts` files, all "imports" paths must be relative
and cannot use a link defined in the `importmap` or `tsconfig.json`.

If a module imported by the workers `.ts` files itself import another module,
the latter should also have a relative path, or else a silencious error will occur.

**Which is very painful to debug**.

Exception can be made for *import type*.
