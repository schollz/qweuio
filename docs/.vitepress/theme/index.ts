import DefaultTheme from 'vitepress/theme'
import CustomLayout from './Layout.vue'
import './style.css'

export default {
  ...DefaultTheme,
  Layout: CustomLayout,
  enhanceApp({ app, router, siteData }) {
    if (typeof window !== 'undefined') {
      const hideNumbersForSingleLine = () => {
        setTimeout(() => {
          const codeBlocks = document.querySelectorAll('div[class*="language-"]')
          codeBlocks.forEach(block => {
            const shiki = block.querySelector('.shiki')
            const lineNumbers = block.querySelector('.line-numbers-wrapper')
            if (shiki && lineNumbers) {
              const lines = shiki.querySelectorAll('.line')
              if (lines.length <= 1) {
                lineNumbers.style.display = 'none'
                shiki.style.left = '-36px'
              }
            }
          })
        }, 100)
      }

      router.onAfterRouteChanged = hideNumbersForSingleLine
      
      // Also run on initial load
      if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', hideNumbersForSingleLine)
      } else {
        hideNumbersForSingleLine()
      }
    }
  }
}
