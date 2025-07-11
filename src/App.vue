<template>
    <div v-focus-visible>
        <header>
            <darkmode-switch />
        </header>
        <main>
            <search-bar v-model="users" :action="() => submit()">
                <advanced-options v-model="advancedOptions" />
            </search-bar>

            <section id="film-results" aria-live="polite" :aria-busy="loading">
                <!-- <plinko :movies="testMovies"></plinko> -->
                <loading-bar v-if="loading" />
                <div v-else-if="submitted">
                    <not-found v-if="notFound" :status="notFoundStatus" />
                    <div v-else>
                        <plinko v-if="movies && movies.length && users" :list="users" :movies="movies" :advancedOptions="advancedOptions"></plinko>
                    </div>

                    <!--
                    <film-result
                        v-else
                        :title="title"
                        :url="url"
                        :imgUrl="imgUrl"
                    /> -->
                </div>
            </section>
        </main>
    </div>
</template>

<script>
import Logo from "./components/Logo.vue";
import DarkmodeSwitch from "./components/DarkmodeSwitch.vue";
import AboutText from "./components/AboutText.vue";
import SearchBar from "./components/SearchBar.vue";
import AdvancedOptions from "./components/AdvancedOptions.vue";
import LoadingBar from "./components/LoadingBar.vue";
import NotFound from "./components/NotFound.vue";
import FilmResult from "./components/FilmResult.vue";
import GoodbyteFooter from "./components/GoodbyteFooter.vue";
import Plinko from "./components/Plinko.vue";

import cache from "./cache.json"

export default {
    name: "App",
    components: {
        Logo,
        DarkmodeSwitch,
        AboutText,
        SearchBar,
        AdvancedOptions,
        LoadingBar,
        NotFound,
        FilmResult,
        GoodbyteFooter,
        Plinko,
    },
    data() {
        return {
            users: "", // list of users to search
            advancedOptions: {}, // advanced option settings

            movies: [
                {
                    film_name: "The Shawshank Redemption",
                    slug: "https://www.imdb.com/title/tt0111161/",
                    image_url:
                        "https://a.ltrbxd.com/resized/film-poster/1/0/4/5/1/1/5/1045115-didi--0-125-0-187-crop.jpg?v=e65b2443a7",
                },
                {
                    film_name: "The Godfather",
                    slug: "https://www.imdb.com/title/tt0068646/",
                    image_url:
                        "https://a.ltrbxd.com/resized/sm/upload/7s/h0/6e/0f/hhsYCYVPy1V0eTjkGNGdvpDO2qk-0-125-0-187-crop.jpg?v=f8d5b0fe80",
                },
                {
                    film_name: "The Dark Knight",
                    slug: "https://www.imdb.com/title/tt0468569/",
                    original_index: 1,
                    image_url:
                        "https://a.ltrbxd.com/resized/film-poster/5/0/8/2/4/6/508246-another-round-0-125-0-187-crop.jpg?v=6227862b11",
                },
                {
                    film_name: "The Shawshank Redemption",
                    slug: "https://www.imdb.com/title/tt0111161/",
                    image_url:
                        "https://a.ltrbxd.com/resized/sm/upload/28/um/1t/jq/dYvyF1RlNokAd1N7Nek0vDpYsV6-0-125-0-187-crop.jpg?v=fc5d71c744",
                },
                {
                    film_name: "The Godfather",
                    slug: "https://www.imdb.com/title/tt0068646/",
                    original_index: 0,
                    image_url:
                        "https://a.ltrbxd.com/resized/film-poster/1/7/9/8/8/8/179888-gueros-0-125-0-187-crop.jpg?v=363c5f44ee",
                },
                {
                    film_name: "The Dark Knight",
                    slug: "https://www.imdb.com/title/tt0468569/",
                    image_url:
                        "https://a.ltrbxd.com/resized/sm/upload/ji/5q/0k/rv/v6xrz4fr92KY1oNC3HsEvrsvR1n-0-230-0-345-crop.jpg?v=973d70bb0c",
                },
                {
                    film_name: "The Shawshank Redemption",
                    slug: "https://www.imdb.com/title/tt0111161/",
                    image_url:
                        "https://a.ltrbxd.com/resized/sm/upload/ji/5q/0k/rv/v6xrz4fr92KY1oNC3HsEvrsvR1n-0-230-0-345-crop.jpg?v=973d70bb0c",
                },
                {
                    film_name: "The Godfather",
                    slug: "https://www.imdb.com/title/tt0068646/",
                    image_url:
                        "https://a.ltrbxd.com/resized/sm/upload/ji/5q/0k/rv/v6xrz4fr92KY1oNC3HsEvrsvR1n-0-230-0-345-crop.jpg?v=973d70bb0c",
                },
                {
                    film_name: "The Dark Knight",
                    slug: "https://www.imdb.com/title/tt0468569/",
                    original_index: 2,
                    image_url:
                        "https://a.ltrbxd.com/resized/sm/upload/ji/5q/0k/rv/v6xrz4fr92KY1oNC3HsEvrsvR1n-0-230-0-345-crop.jpg?v=973d70bb0c",
                },
                {
                    film_name: "The Godfather",
                    slug: "https://www.imdb.com/title/tt0068646/",
                    image_url:
                        "https://a.ltrbxd.com/resized/sm/upload/ji/5q/0k/rv/v6xrz4fr92KY1oNC3HsEvrsvR1n-0-230-0-345-crop.jpg?v=973d70bb0c",
                }
            ],

            info: "", // json blob from AJAX request
            currentHash: null, // track most recent request

            loading: false, // loading statehood
            submitted: true, // submit statehood

            notFound: false, // not found error
            emptyintersect: false, // intersect error
            timeout: false, // timeout error
        };
    },
    created() {
        // get starting values from url params (if exist)
        const urlParams = new URLSearchParams(window.location.search);

        const users = urlParams.getAll("u");
        const intersect = urlParams.get("i");
        const ignore = urlParams.get("ignore");

        if (users.length > 0) {
            this.users = users.toString();
        }

        if (intersect != null) {
            this.advancedOptions = { selectionMode: "Intersect" };
        }

        if (ignore != null) {
            let ignoreList = ignore.split(",");

            if (
                ignoreList.includes("unreleased") ||
                ignoreList.includes("true")
            ) {
                this.advancedOptions["unreleased"] = false;
            }

            if (ignoreList.includes("shorts")) {
                this.advancedOptions["shortFilms"] = false;
            }

            if (ignoreList.includes("feature")) {
                this.advancedOptions["featureLength"] = false;
            }
        }

        this.submit();
    },
    methods: {
        /**
         * main function, runs process to get random film
         */
        submit() {
            this.notFound = false;

            if (this.users == "") {
                // reset state if form submitted with empty input field
                this.reset();
                return;
            }

            let userlist = this.users.split(/(?:,| )+/); // split input field on space or comma

            if (userlist.length < 1) {
                // second check for non empty input field (only whotespace or commas entered)
                this.submitted = false;
                return;
            }

            document.body.classList.remove("done");
            document.body.classList.add("entered");

            this.submitted = true;
            this.loading = true;

            let ignoreList = [];
            if (this.advancedOptions["unreleased"] == false) {
                ignoreList.push("unreleased");
            }
            if (this.advancedOptions["shortFilms"] == false) {
                ignoreList.push("shorts");
            }
            if (this.advancedOptions["featureLength"] == false) {
                ignoreList.push("feature");
            }

            let apiUrl = "/api?users=" + userlist.join("&users=")
            let clientUrl = "?u=" + userlist.join("&u=");

            if (ignoreList.length > 0) {
                apiUrl += "&ignore=" + ignoreList.join(",");
                clientUrl += "&ignore=" + ignoreList.join(",");
            }

            if (this.advancedOptions["selectionMode"] == "Intersect") {
                apiUrl += "&intersect=true";
                clientUrl += "&i=true";
            }

            window.history.replaceState(null, null, clientUrl);

            try {
                let vue = this;
                let hash = this.hashCode(apiUrl);
                console.log("url: " + apiUrl + "\nhash: " + hash);

                // If localhost, use dummy data on this component
                if (window.location.hostname == "localhost") {
                    this.loading = false;
                    this.movies = [...cache, ...cache]
                    return;
                }

                vue.currentHash = hash;
                fetch(apiUrl)
                    .then(function (res) {
                        // check if new request has been sent since submitting
                        if (vue.currentHash != hash) {
                            return "ignoreOldRequest";
                        }

                        document.body.classList.remove("entered");
                        document.body.classList.add("done");

                        // bad response
                        if (res.status != 200) {
                            if (res.status == 406) {
                                vue.emptyintersect = true;
                            }

                            if (res.status == 502) {
                                vue.timeout = true;
                            }

                            vue.notFound = true;
                            vue.loading = false;

                            return "";
                        }

                        // good response

                        setTimeout(function () {
                            vue.loading = false;
                        }, 200); // wait 200ms to load text to allow for image to preload

                        return res.json();
                    })
                    .then(function (json) {
                        if (json == "ignoreOldRequest") {
                            return;
                        }

                        if (!vue.notfound) {
                            // preload image
                            var pre_image = new Image();
                            pre_image.src = json.image_url;
                        }

                        vue.movies = json.films;
                    });
            } catch (e) {
                alert(
                    "Something went wrong. Please try again in a moment. Error:" +
                    e,
                    "An error occured",
                );
            }
        },

        /**
         * reset the state
         */
        reset() {
            window.history.replaceState(null, null, "/");
            this.loading = false;
            this.submitted = false;
            document.body.classList.remove("done");
            document.body.classList.remove("entered");
        },

        /**
         * hashing function to help manage requests
         */
        hashCode(s) {
            let h;
            for (let i = 0; i < s.length; i++)
                h = (Math.imul(31, h) + s.charCodeAt(i)) | 0;
            return h + Date.now();
        },
    },
    computed: {
        notFoundStatus: function () {
            if (this.emptyintersect) {
                return "no-intersect";
            } else if (this.timeout) {
                return "timeout";
            } else if (
                !(
                    this.advancedOptions["unreleased"] &&
                    this.advancedOptions["shortFilms"] &&
                    this.advancedOptions["featureLength"]
                )
            ) {
                return "possibly-ignored";
            } else {
                return "other";
            }
        },
    },
};
</script>

<style>
:root {
    --background: #fff;
    --foreground: #2c3e50;

    --primary: oklch(43.48% 0.17 260.2);
    --secondary: #000;
    --tertiary: oklch(62.39% 0.181 258.33);

    --white: #fff;
    --off-white: #ebebeb;
    --black: #000;
}

.dark {
    --background: oklch(14.38% 0.007 256.88);
    --foreground: oklch(99.4% 0 0);

    --secondary: #fff;
    --tertiary: oklch(62.39% 0.181 258.33);
}

body {
    background: var(--background);

    margin: 30px 0 0 0;

    font-family:
        Avenir,
        -apple-system,
        Helvetica,
        Arial,
        sans-serif;
    color: var(--foreground);
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;

    transition:
        color ease-in-out 0.25s,
        background-color ease-in-out 0.25s;
}

main {
    min-height: calc(100vh - 100px);
    transform: translateY(40px);
    transition: transform 1.2s cubic-bezier(0.82, 0.01, 0.45, 1);
}

.entered main,
.done main {
    transform: none;
}

#film-results {
    margin: -20px 0;
    transform: translateY(-160px);
    transition: transform 0.3s ease;
}

#film-results.advanced-active {
    background-color: var(--background);
    position: relative;
    z-index: 99;

    transform: translateY(0px);
    transition: transform 0.3s ease;
}

::selection {
    background: var(--primary);
    color: var(--white);
}

@media (prefers-reduced-motion) {
    * {
        animation: none !important;
        transition: none !important;
    }
}
</style>
