#! /bin/bash
# Codemod to minimize bundle size
find src -name '*.tsx' -print | xargs jscodeshift -t node_modules/@material-ui/codemod/lib/v4.0.0/top-level-imports.js --parser=tsx
# find src -name '*.ts' -print | xargs jscodeshift -t node_modules/@material-ui/codemod/lib/v4.0.0/top-level-imports.js --parser=ts

find src -name '*.tsx' -print | xargs jscodeshift -t node_modules/@material-ui/codemod/lib/v4.0.0/optimal-imports.js --parser=tsx
