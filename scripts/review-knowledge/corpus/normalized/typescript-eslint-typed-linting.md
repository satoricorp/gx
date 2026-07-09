# Linting with Type Information

Some typescript-eslint rules utilize TypeScript's type checking APIs to provide much deeper insights into your code. This requires TypeScript to analyze your entire project instead of just the file being linted. As a result, these rules are slower than traditional lint rules but are much more powerful.

To enable typed linting, there are two small changes you need to make to your config file:

- Flat Config
- Legacy Config

- Add `TypeChecked`to the name of any preset configs you're using, namely`recommended`,`strict`, and`stylistic`.
- Add `languageOptions.parserOptions`to tell our parser how to find the TSConfig for each source file.

```
import js from '@eslint/js';
import { defineConfig } from 'eslint/config';
import tseslint from 'typescript-eslint';
export default defineConfig({
 files: ['**/*.{js,ts}'],
 extends: [
 js.configs.recommended,
 tseslint.configs.recommended,
 tseslint.configs.recommendedTypeChecked,
 ],
 languageOptions: {
 parserOptions: {
 projectService: true,
 },
 },
});
```
In more detail:

- `tseslint.configs.recommendedTypeChecked`is a shared configuration. It contains recommended rules that additionally require type information.
- `parserOptions.projectService: true`indicates to ask TypeScript's type checking service for each source file's type information (see Parser > Project Service).

- Add `-type-checked`to the name of any preset configs you're using, namely`recommended`,`strict`, and`stylistic`.
- Add `parserOptions`to tell our parser how to find the TSConfig for each source file.

```
/* eslint-env node */
module.exports = {
 extends: [
 'eslint:recommended',
 'plugin:@typescript-eslint/recommended',
 'plugin:@typescript-eslint/recommended-type-checked',
 ],
 plugins: ['@typescript-eslint'],
 parser: '@typescript-eslint/parser',
 parserOptions: {
 projectService: true,
 tsconfigRootDir: __dirname,
 },
 root: true,
};
```
In more detail:

- `plugin:@typescript-eslint/recommended-type-checked`is a shared configuration. It contains recommended rules that additionally require type information.
- `parserOptions.projectService: true`indicates to ask TypeScript's type checking service for each source file's type information (see Parser >- `projectService`).
- `parserOptions.tsconfigRootDir`tells our parser the absolute path of your project's root directory (see Parser >- `tsconfigRootDir`).

Your ESLint config file may start receiving a parsing error about type information. See our TSConfig inclusion FAQ.

With that done, run the same lint command you ran before. You may see new rules reporting errors based on type information!

## Shared Configurations

If you enabled the `strict` shared config and/or `stylistic` shared config in a previous step, be sure to replace them with `strictTypeChecked` and `stylisticTypeChecked` respectively to add their type-checked rules.

- Flat Config
- Legacy Config

```
export default defineConfig({
 files: ['**/*.{js,ts}'],
 extends: [
 js.configs.recommended,
 tseslint.configs.strict,
 tseslint.configs.stylistic,
 tseslint.configs.strictTypeChecked,
 tseslint.configs.stylisticTypeChecked,
 ],
 // ...
});
```
```
/* eslint-env node */
module.exports = {
 extends: [
 'eslint:recommended',
 'plugin:@typescript-eslint/strict',
 'plugin:@typescript-eslint/stylistic',
 'plugin:@typescript-eslint/strict-type-checked',
 'plugin:@typescript-eslint/stylistic-type-checked',
 ],
 // ...
};
```
You can read more about the rules provided by typescript-eslint in our rules docs and shared configurations docs.

## Performance

Typed rules come with a catch. By using typed linting in your config, you incur the performance penalty of asking TypeScript to do a build of your project before ESLint can do its linting. For small projects this takes a negligible amount of time (a few seconds or less); for large projects, it can take longer.

Most of our users do not mind this cost as the power and safety of type-aware static analysis rules is worth the tradeoff. Additionally, most users primarily consume lint errors via IDE plugins which, through caching, do not suffer the same penalties. This means that generally they usually only run a complete lint before a push, or via their CI, where the extra time often doesn't matter.

**We strongly recommend you do use type-aware linting**, but the above information is included so that you can make your own, informed decision.
See Troubleshooting > Typed Linting > Performance for more information.

## Troubleshooting

If you're having problems with typed linting, please see our Troubleshooting FAQs and in particular Troubleshooting > Typed Linting.

For details on the parser options that enable typed linting, see:

- Parser > `projectService`: our recommended option, with settings to customize TypeScript project information
- Parser > `project`: an older option that can be used as an alternative
