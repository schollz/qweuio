<script setup>
import DefaultTheme from 'vitepress/theme'
import { ref, onMounted } from 'vue'

const DefaultLayout = DefaultTheme.Layout

const menuOpen = ref(false)
const isDark = ref(false)

onMounted(() => {
    isDark.value = document.documentElement.classList.contains('dark')
})

const closeMenu = () => {
    menuOpen.value = false
}
</script>

<template>
    <DefaultLayout>

        <template #layout-bottom>
            <footer class="text-center mt-2 mb-2">
                <p>Built with Brutalism & VitePress</p>
            </footer>
        </template>
    </DefaultLayout>
</template>
