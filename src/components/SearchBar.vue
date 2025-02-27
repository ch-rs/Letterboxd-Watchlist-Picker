<template>
    <section>
        <label for="userbox"> Username(s): </label>
        <div class="options" v-if="shortcuts && shortcuts.length > 0">
            <button v-for="shortcut in shortcuts" :key="shortcut" v-on:click="updateValue(shortcut)">
                {{ shortcut }}
            </button>
        </div>

        <div class="form-container">
            <input id="userbox" class="userfield" type="text" placeholder="ex: holopollock, qjack" :value="value"
                v-on:keyup.enter="action()" v-on:input="updateValue($event.target.value)" />
            <button v-on:click="action()">Submit</button>
        </div>

        <slot />
    </section>
</template>

<script>
import ding from '../utils/ding';

export default {
    name: "SearchBar",
    props: ["value", "action"],
    data() {
        return {
            shortcuts: []
        };
    },
    created() {
        this.loadShortcuts();
        
        // Add event listener for shortcut updates
        window.addEventListener('shortcuts-updated', this.handleShortcutsUpdated);
    },
    beforeDestroy() {
        // Clean up event listener when component is destroyed
        window.removeEventListener('shortcuts-updated', this.handleShortcutsUpdated);
    },
    methods: {
        updateValue: function (value) {
            this.$emit("input", value);
        },
        setValue: function (value) {
            this.value = value;
        },
        beep: function () {
            ding(Math.max(Math.random(), 0.3), Math.random() * 500);
        },
        loadShortcuts: function() {
            const savedShortcuts = localStorage.getItem('userShortcuts');
            if (savedShortcuts) {
                try {
                    this.shortcuts = JSON.parse(savedShortcuts);
                } catch (e) {
                    console.error('Failed to parse shortcuts from localStorage', e);
                    this.shortcuts = [];
                }
            } else {
                // Default shortcuts if none are saved
                this.shortcuts = [
                ];
            }
        },
        handleShortcutsUpdated(event) {
            this.shortcuts = event.detail;
        }
    },
};
</script>

<style scoped>
label {
    visibility: hidden;
    display: block;
}

.options {
    display: flex;
    flex-direction: column;
    margin: 0 auto 2rem;
    gap: 1em;
    max-width: 300px;
}

.options button {
    border-radius: 4px !important;
}

.form-container {
    display: flex;
    align-items: center;
    justify-content: center;
}

.userfield {
    font-family: Avenir, Helvetica, Arial, sans-serif;
    font-size: 1.1rem;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    line-height: 2.8rem;

    background: var(--off-white);
    padding: 0 1rem;
    min-width: 220px;
    border: 0;
    border-radius: 4px 0 0 4px;
    outline: none;
}

.userfield:active,
[v-focus-visible="true"] .userfield:focus,
[v-focus-visible="true"] .userfield:focus-within {
    box-shadow: inset 0 0 0 3px var(--primary);
}

::placeholder {
    opacity: 0.6;
}

button {
    color: var(--white);
    font-family: Avenir, Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    font-size: 0.8rem;
    font-weight: 900;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    line-height: 2.8rem;

    display: inline-block;
    cursor: pointer;
    padding: 0 1rem;
    border: 0;
    border-radius: 0 4px 4px 0;
    outline: none;

    background: var(--secondary);
    transition: background-color ease-in-out 0.2s;
    /* darkmode transition */
}

.dark button {
    background: var(--tertiary);
}

button:hover,
[v-focus-visible="true"] button:focus {
    background: var(--primary);
    transition: none;
}

@media screen and (max-width: 360px) {
    .userfield {
        min-width: 0px;
        width: 50%;
    }
}
</style>
