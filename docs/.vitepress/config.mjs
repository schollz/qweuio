import { defineConfig } from 'vitepress'
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const grammar = JSON.parse(
  fs.readFileSync(path.resolve(__dirname, './theme/tli.tmLanguage.json'), 'utf-8')
)

export default defineConfig({
  title: 'asdfgh',
  description: 'Documentation',
  themeConfig: {
    appearance: true,
    sidebar: false,
    nav: [
      { text: 'tli', link: '#tli' },
      { text: 'Features', link: '#features' },
      { text: 'Contact', link: '#contact' }
    ]
  },
  markdown: {
    code: {
      highlight: async (code, lang) => {
        const { getHighlighter } = await import('shiki')

        const highlighter = await getHighlighter({
          themes: ['github-dark'],
          langs: [
            // Register TLI
            {
              id: 'tli',
              scopeName: 'source.tli',
              grammar
            }
          ]
        })

        // 🧪 DEBUG (optional)
        console.log('[SHIKI] Loaded languages:', highlighter.getLoadedLanguages())

        return highlighter.codeToHtml(code, {
          lang,
          theme: 'github-dark'
        })
      }
    }
  }
})
