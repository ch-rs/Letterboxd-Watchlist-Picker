<template>
	<section id="advanced-section" aria-haspopup="true" aria-expanded="false">
		<button id="show-advanced" v-on:click="activateAdvanced()">
			Advanced Options
		</button>
		<div id="collapsable">
			<div>
				<input type="checkbox" id="relative" value="Relative" name="relative" v-model="relative" v-on:change="updateValues()" checked />
				<label for="relative" class="relative-label">
					Prioritise older movies
				</label>
			</div>

			<br />

			<input
				type="radio"
				id="union"
				value="Union"
				name="union-intersect"
				tabindex="-1"
				v-model="selectionMode"
				v-on:change="updateValues()"
				selected
			/>
			<label for="union" class="union-intersection-label">
				<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 131 87">
					<g stroke="#000" stroke-width="3">
						<g transform="translate(-796 -208)">
							<circle cx="839.7" cy="251.7" r="43.7" stroke="none"/>
							<circle cx="839.7" cy="251.7" r="42.2" fill="none"/>
						</g>
						<g transform="translate(-796 -208)">
							<circle cx="883.3" cy="251.7" r="43.7" stroke="none"/>
							<circle cx="883.3" cy="251.7" r="42.2" fill="none"/>
						</g>
						<g fill="none" transform="translate(-796 -208)">
							<circle cx="839.7" cy="251.7" r="43.7" stroke="none"/>
							<circle cx="839.7" cy="251.7" r="42.2"/>
						</g>
						<path fill="#fff" d="M15 13l73 73zM2 40l46 46zm35 45L2 50zM20 8l77 77zm56 76L10 18zM3 31l53 53zM27 5l77 77zM6 24l57 57zm16 56L8 66zM35 3l76 76zm9-1l72 72zm11 1l66 66zm13 3l57 57zm7-3l53 53zm9-1l45 45zm10 0l35 35zm16 6l13 13z"/>
					</g>
				</svg>
				Union
			</label>
			<input
				type="radio"
				id="intersect"
				value="Intersect"
				name="union-intersect"
				tabindex="-1"
				v-model="selectionMode"
				v-on:change="updateValues()"
			/>
			<label for="intersect" class="union-intersection-label">
				<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 131 87">
					<g stroke="#000" stroke-width="3">
						<g transform="translate(-786 -52)">
							<circle cx="829.7" cy="95.7" r="43.7" stroke="none"/>
							<circle cx="829.7" cy="95.7" r="42.2" fill="none"/>
						</g>
						<g transform="translate(-786 -52)">
							<circle cx="873.3" cy="95.7" r="43.7" stroke="none"/>
							<circle cx="873.3" cy="95.7" r="42.2" fill="none"/>
						</g>
						<g fill="none" transform="translate(-786 -52)">
							<circle cx="829.7" cy="95.7" r="43.7" stroke="none"/>
							<circle cx="829.7" cy="95.7" r="42.2"/>
						</g>
						<path fill="#fff" d="M47 55l22 22zm-2-12l29 29zm1-9l33 33zm3-7l33 33zm3-7l33 33zm5-5l29 29zm5-5l22 22z"/>
					</g>
				</svg>
				Intersect
			</label>

			<div class="checkbox-container">
				<h4>Include in results:</h4>
				
				<div>
					<input
						type="checkbox"
						id="ignore-unreleased"
						tabindex="-1"
						v-model="unreleased"
						v-on:change="updateValues()"
						v-on:keyup.enter="unreleased = !unreleased; updateValues()"
						selected
					/>
					<label for="ignore-unreleased">
						Unreleased films
					</label>
				</div>
				<div>
					<input
						type="checkbox"
						id="ignore-short"
						tabindex="-1"
						v-model="shortFilms"
						v-on:change="updateValues()"
						v-on:keyup.enter="shortFilms = !shortFilms; updateValues()"
						selected
					/>
					<label for="ignore-short">
						Short films
					</label>
				</div>
				<div>
					<input
						type="checkbox"
						id="ignore-feature"
						tabindex="-1"
						v-model="featureLength"
						v-on:change="updateValues()"
						v-on:keyup.enter="featureLength = !featureLength; updateValues()"
						selected
					/>
					<label for="ignore-feature">
						Feature-length films
					</label>
				</div>
			</div>

			<!-- New Shortcuts Section -->
			<div class="shortcuts-container">
				<h4>Shortcuts:</h4>
				<div class="shortcut-input">
					<input 
						type="text" 
						id="new-shortcut" 
						v-model="newShortcut" 
						placeholder="Enter username/path"
						:tabindex="advancedOpen ? 0 : -1"
					/>
					<button 
						@click="addShortcut" 
						:tabindex="advancedOpen ? 0 : -1"
					>
						Add
					</button>
				</div>
				<div class="shortcuts-list">
					<div v-for="(shortcut, index) in shortcuts" :key="index" class="shortcut-item">
						<span>{{ shortcut }}</span>
						<button 
							@click="removeShortcut(index)" 
							class="remove-shortcut"
							:tabindex="advancedOpen ? 0 : -1"
						>
							✕
						</button>
					</div>
					<div v-if="shortcuts.length === 0" class="no-shortcuts">
						No shortcuts added yet
					</div>
				</div>
			</div>
		</div>
	</section>
</template>


<script>
	export default
	{
		name: 'AdvancedOptions',
		props: [ 'value' ],
		data()
		{
			return {
				relative: true,
				selectionMode: "Union",
				unreleased: true,
				shortFilms: true,
				featureLength: true,
				advancedOpen: false,
				shortcuts: [],
				newShortcut: ""
			};
		},
		created()
		{
			if (this.value)
			{
				if ('selectionMode' in this.value)
				{
					this.selectionMode = this.value['selectionMode'];
				}

				if ('unreleased' in this.value)
				{
					this.unreleased = this.value['unreleased'];
				}

				if ('shortFilms' in this.value)
				{
					this.shortFilms = this.value['shortFilms'];
				}

				if ('featureLength' in this.value)
				{
					this.featureLength = this.value['featureLength'];
				}

				if ('relative' in this.value)
				{
					this.relative = this.value['relative'];
				}
			}

			this.loadShortcuts();
			this.updateValues();
		},
		methods:
		{
			activateAdvanced()
			{
				let collapsable = document.getElementById("collapsable");
				let results = document.getElementById("film-results");

				collapsable.classList.toggle("active");
				results.classList.toggle("advanced-active");

				if (this.advancedOpen)
				{
					document.getElementById("union").tabIndex = -1;
					document.getElementById("intersect").tabIndex = -1;
					document.getElementById("ignore-unreleased").tabIndex = -1;
					document.getElementById("ignore-short").tabIndex = -1;
					document.getElementById("ignore-feature").tabIndex = -1;
					document.getElementById("relative").tabIndex = -1;
					document.getElementById("advanced-section").setAttribute("ariaexpanded", "false");
					this.advancedOpen = false;
				}
				else
				{
					document.getElementById("union").tabIndex = 0; // allow tab-to-focus
					document.getElementById("intersect").tabIndex = 0;
					document.getElementById("ignore-unreleased").tabIndex = 0;
					document.getElementById("ignore-short").tabIndex = 0;
					document.getElementById("ignore-feature").tabIndex = 0;
					document.getElementById("relative").tabIndex = 0;
					document.getElementById("advanced-section").setAttribute("ariaexpanded", "true");
					this.advancedOpen = true;
				}
			},

			updateValues()
			{
				this.$emit('input', {
					'selectionMode': this.selectionMode,
					'relative': this.relative,
					'unreleased': this.unreleased,
					'shortFilms': this.shortFilms,
					'featureLength': this.featureLength
				});
			},

			// New methods for shortcuts
			loadShortcuts() {
				const savedShortcuts = localStorage.getItem('userShortcuts');
				if (savedShortcuts) {
					try {
						this.shortcuts = JSON.parse(savedShortcuts);
					} catch (e) {
						console.error('Failed to parse shortcuts from localStorage', e);
						this.shortcuts = [];
					}
				}
			},

			saveShortcuts() {
				localStorage.setItem('userShortcuts', JSON.stringify(this.shortcuts));
				// Dispatch a custom event that other components can listen for
				window.dispatchEvent(new CustomEvent('shortcuts-updated', {
					detail: this.shortcuts
				}));
			},

			addShortcut() {
				if (this.newShortcut.trim()) {
					// Prevent duplicates
					if (!this.shortcuts.includes(this.newShortcut.trim())) {
						this.shortcuts.push(this.newShortcut.trim());
						this.saveShortcuts();
					}
					this.newShortcut = "";
				}
			},

			removeShortcut(index) {
				this.shortcuts.splice(index, 1);
				this.saveShortcuts();
			}
		}
	}
</script>

<style scoped>
	#show-advanced
	{
		display: block;
		position: relative;
		margin: auto;
		padding: 0 1rem;
		border: 0;
		border-radius: 0 4px 4px 0;
		background: none;
		outline: none;
		cursor: pointer;

		color: var(--foreground);
		font-family: Avenir, -apple-system, Helvetica, Arial, sans-serif;
		-webkit-font-smoothing: antialiased;
		-moz-osx-font-smoothing: grayscale;
		font-size: 0.8rem;
		font-weight: 400;
		letter-spacing: 0.04em;
		line-height: 2.8rem;

		background-size: calc(100% - 2rem) 1px;
		background-position: 1rem 2.3em;
		background-repeat: no-repeat;

		opacity: 0.6;
		transform: translateY(-0.12rem); /* adjustment */
		transition: color ease-in-out 0.25s;
	}

	#show-advanced:hover,
	[v-focus-visible=true] #show-advanced:focus
	{
		opacity: 1;
		color: var(--primary);
		background-image: linear-gradient(var(--primary), var(--primary));
		transition: none;
	}

	[v-focus-visible=true] #show-advanced:focus::after
	{
		content: "";
		display: block;
		position: absolute;
		top: 3px;
		left: -3px;
		right: -3px;
		bottom: 3px;

		border: 2px solid var(--primary);
		border-radius: 6px;
	}

	#collapsable
	{
		opacity: 0;
		transition: 0.3s ease;
		transform: translateY(0px);
	}

	#collapsable.active
	{
		opacity: 1;
		transform: translateY(5px);
	}


	/*==========================================================================
		options
	==========================================================================*/

	.union-intersection-label
	{
		display: inline;
		cursor: pointer;
		margin: 4px;
		padding: 10px;
		border: 2px solid transparent;
		border-radius: 4px;

		color: var(--black);
		font-weight: bold;

		background: var(--tertiary);
	}

	.union-intersection-label svg
	{
		width: 40px;
		transform: translateY(8.5px);
		margin-right: 4px;
		fill: transparent;
	}

	#union, #intersect
	{
		opacity: 0;
		width: 0px;
		margin: 0;
	}

	input:checked + .union-intersection-label
	{
		color: var(--white);
		background: var(--primary);
		border-color: var(--secondary);
	}

	#collapsable input:checked + .union-intersection-label svg
	{
		fill: var(--white);
	}

	[v-focus-visible=true] input:focus + .union-intersection-label
	{
		box-shadow: 0 0 0 3px var(--background), 0 0 0 5px var(--primary);
	}

	.checkbox-container
	{
		padding: 1rem 0.8rem;
	}

	.checkbox-container h4
	{
		margin-bottom: 0.5rem;
	}

	.checkbox-container div
	{
		display: inline-block;
		background-color: var(--off-white);
		color: var(--black);
		margin: 5px 5px;
		padding: 5px 10px;
		border-radius: 6px;
	}

	.dark .checkbox-container div
	{
		opacity: 0.6;
	}

	[v-focus-visible=true] .checkbox-container div:focus-within
	{
		box-shadow: 0 0 0 3px var(--background), 0 0 0 5px var(--primary);
		opacity: 1;
	}

	.checkbox-container input
	{
		margin-right: 10px;
		cursor: pointer;
	}

	.checkbox-container label
	{
		display: inline-block;
		background: none;
		text-align: left;
		cursor: pointer;
	}

	.checkbox-container span
	{
		font-size: 12px;
		display: inline-block;
	}

	/* New styles for shortcuts */
	.shortcuts-container {
		margin-top: 1rem;
		padding: 1rem 0.8rem;
	}

	.shortcuts-container h4 {
		margin-bottom: 0.5rem;
	}

	.shortcut-input {
		display: flex;
		margin-bottom: 0.8rem;
	}

	.shortcut-input input {
		flex: 1;
		padding: 5px 10px;
		border: 1px solid var(--tertiary);
		border-radius: 4px 0 0 4px;
		font-family: inherit;
	}

	.shortcut-input button {
		padding: 5px 10px;
		background-color: var(--primary);
		color: var(--white);
		border: none;
		border-radius: 0 4px 4px 0;
		cursor: pointer;
	}

	.shortcuts-list {
		max-height: 150px;
		overflow-y: auto;
		border: 1px solid var(--tertiary);
		border-radius: 4px;
		background-color: var(--off-white);
	}

	.shortcut-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 5px 10px;
		border-bottom: 1px solid var(--tertiary);
	}

	.shortcut-item:last-child {
		border-bottom: none;
	}

	.remove-shortcut {
		background: none;
		border: none;
		color: var(--black);
		cursor: pointer;
		font-size: 0.8rem;
		padding: 2px 5px;
		border-radius: 50%;
	}

	.remove-shortcut:hover {
		background-color: var(--tertiary);
	}

	.no-shortcuts {
		padding: 10px;
		text-align: center;
		color: var(--black);
		opacity: 0.6;
	}

	.dark .shortcuts-list {
		background-color: var(--background);
		color: var(--foreground);
	}

	.dark .shortcut-item {
		border-color: var(--background);
	}

	.dark .remove-shortcut {
		color: var(--foreground);
	}

	.dark .remove-shortcut:hover {
		background-color: var(--background);
	}
</style>
