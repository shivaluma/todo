import { dirname } from "path"
import { fileURLToPath } from "url"
import { FlatCompat } from "@eslint/eslintrc"

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

const compat = new FlatCompat({
  baseDirectory: __dirname,
})

const eslintConfig = [
  ...compat.extends("next/core-web-vitals", "next/typescript"),
  ...compat.extends("plugin:prettier/recommended"),
  ...compat.extends("plugin:jsx-a11y/recommended"),
  ...compat.extends(
    "plugin:react/recommended",
    "plugin:react-hooks/recommended"
  ),
  {
    rules: {
      "react/react-in-jsx-scope": "off",
    },
  },
]

export default eslintConfig
