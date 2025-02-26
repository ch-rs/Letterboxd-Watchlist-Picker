<template>
    <div>
        <button class="drop" @click="dropBall">Drop!</button>
        <div class="parent" ref="myCanvas">
            <div v-if="movies && widthPercentages" class="images">
                <div v-for="(movie, i) in movies" :key="i" :class="{
                    'image': true,
                    'active': activeMovieIndex === i
                }" :style="{ width: widthPercentages[i] + '%' }">
                    <a :href="movie.slug" target="_blank">
                        <img :src="movie.image_url" />
                    </a>
                    <span class="movie-name">
                        <span v-text="movie.film_name"></span>
                        <a :href="movie.slug" @click="$event.target.classList.add('clicked')" target="_blank">++++</a>
                        <a :href="thisList" @click="$event.target.classList.add('clicked')" target="_blank">-------</a>
                    </span>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import ding from '../utils/ding';
import distinctColors from 'distinct-colors'


const random = Matter.Common.random;
const Engine = Matter.Engine;
const World = Matter.World;
const Bodies = Matter.Bodies;

let sketch = function (p, parent) {
    let engine,
        world,
        particles = [],
        plinkos = [],
        bounds = [],
        movies = [],
        cols = 11,
        spacing,
        palette,
        frame = 0,
        widthPercentages = [],
        rows = 10,
        rowStart = 100,
        particleSize = 12,
        slotWidth,
        ballPositions = {},
        alerted = false,
        plinkoSize = 14;

    p.setup = function () {
        const c = p.createCanvas(612, 800);
        c.parent(parent);

        engine = Engine.create();
        world = engine.world;
        world.gravity.y = 0.45;

        spacing = p.width / cols;

        Matter.Events.on(engine, "collisionStart", function (event) {
            for (let i = 0; i < event.pairs.length; i++) {
                let pair = event.pairs[i];
                let velocityA = pair.bodyA.velocity;
                let velocityB = pair.bodyB.velocity;

                // Calculate the magnitude of the velocities
                let speedA = Math.sqrt(velocityA.x * velocityA.x + velocityA.y * velocityA.y);
                let speedB = Math.sqrt(velocityB.x * velocityB.x + velocityB.y * velocityB.y);

                // Use the maximum speed of the two bodies for the ding function
                let maxSpeed = Math.max(speedA, speedB);

                let vCeiling = Math.min(maxSpeed / 10, 0.9);
                let v = Math.max(vCeiling, 0.3);

                let frequency = 500;
                if (pair.bodyA.label === 'plinko') {
                    frequency = pair.bodyA.frequency;
                } else if (pair.bodyB.label === 'plinko') {
                    frequency = pair.bodyB.frequency;
                }

                // Call the ding function with the velocity-based volume and frequency

                p.preDraw()

                ding(v, frequency);
            }
        });

        /*
        Matter.Events.on(engine, "afterUpdate", function () {
            for (let i = 0; i < particles.length; i++) {
                // Check if the particle's speed is close to zero
                if (stopped == false && particles[i].body.speed < 0.01) {
                    stopped = true;
                }
            }
        });
        */


    };

    p.populate = function (newMovies) {
        cols = newMovies.length;
        p.movies = newMovies;

        // Calculate the width of each slot
        setTimeout(() => {
            // Check if we have valid original_index values
            const hasValidOriginalIndices = p.movies.some(movie =>
                movie.original_index !== undefined &&
                movie.original_index >= 0 &&
                movie.original_index < 3
            );

            if (hasValidOriginalIndices) {
                // Define width multipliers for the oldest movies
                const topMovieMultipliers = [2.5, 1.8, 1.4]; // Multipliers for movies with OriginalIndex 0, 1, 2

                // Calculate total width units
                let totalWidthUnits = 0;
                for (let i = 0; i < p.movies.length; i++) {
                    const originalIndex = p.movies[i].original_index;
                    if (originalIndex !== undefined && originalIndex >= 0 && originalIndex < 3) {
                        // This is one of the top 3 oldest movies
                        totalWidthUnits += topMovieMultipliers[originalIndex];
                    } else {
                        // Regular movie
                        totalWidthUnits += 1;
                    }
                }

                // Calculate the base unit width
                const baseUnitWidth = p.width / totalWidthUnits;

                // Assign positions and widths to each movie
                let currentX = 0;
                for (let i = 0; i < p.movies.length; i++) {
                    const originalIndex = p.movies[i].original_index;
                    let movieWidth;

                    if (originalIndex !== undefined && originalIndex >= 0 && originalIndex < 3) {
                        // This is one of the top 3 oldest movies
                        movieWidth = baseUnitWidth * topMovieMultipliers[originalIndex];
                    } else {
                        // Regular movie
                        movieWidth = baseUnitWidth;
                    }

                    p.movies[i].x = currentX;
                    p.movies[i].width = movieWidth;

                    // Calculate width as percentage for CSS
                    widthPercentages[i] = (movieWidth / p.width) * 100;

                    currentX += movieWidth;
                }
            } else {
                // No valid original indices - use equal widths for all movies
                const equalWidth = p.width / p.movies.length;

                for (let i = 0; i < p.movies.length; i++) {
                    p.movies[i].x = equalWidth * i;
                    p.movies[i].width = equalWidth;
                    widthPercentages[i] = 100 / p.movies.length;
                }
            }

            slotWidth = p.width / p.movies.length; // Keep this for reference

            p.createPlinkos();
            p.createBoundaries();

            p.preDraw();
        }, 1000);
    };

    p.getWidthPercentages = function () {
        return widthPercentages
    }

    p.newParticle = function () {
        alerted = false
        const part = new Particle(random(5, p.width - 5), 0, particleSize);
        particles.push(part);
    };

    p.createBoundaries = function () {
        let b = new Boundary(p.width / 2, p.height + 50, p.width, 100);
        bounds.push(b);

        // Create boundaries based on movie positions and widths
        for (let i = 0; i <= cols; i++) {
            let x;
            if (i === 0) {
                // First boundary at left edge
                x = 0;
            } else if (i === cols) {
                // Last boundary at right edge
                x = p.width;
            } else {
                // Boundaries between movies
                x = p.movies[i - 1].x + p.movies[i - 1].width;
            }

            const h = 60;
            const w = 3;
            const y = p.height - h / 2;
            b = new Boundary(x, y, w, h);
            bounds.push(b);
        }
    };

    p.createPlinkos = function () {
        let count = 0
        for (let j = 0; j < rows; j++) {
            for (let i = 0; i < cols + 1; i++) {
                count++
            }
        }

        palette = distinctColors({
            count,
            lightMin: 30,
            lightMax: 45,
            chromaMin: 20,
            chromaMax: 100,
            hueMin: 190,
            hueMax: 300
        });

        for (let j = 0; j < rows; j++) {
            for (let i = 0; i < cols + 1; i++) {
                let x = i * spacing;
                if (j % 2 == 0) {
                    x += spacing / 2;
                }
                const y = rowStart + j * spacing;
                const p = new Plinko(x, y, plinkoSize + (Math.random() * 5 - 2.5));
                plinkos.push(p);
            }
        }
    };

    p.checkIfStopped = function (particle) {
        const velocity = particle.body.velocity;
        const speed = Math.sqrt(velocity.x * velocity.x + velocity.y * velocity.y);
        return speed < 0.005; // Threshold for considering the particle as stopped
    }

    p.getSegmentIndex = function (particle) {
        let sw = p.width / cols;
        const x = particle.body.position.x + 5;
        return Math.floor(x / sw);
    }

    p.preDraw = function () {
        p.background('#0a0a0a')

        for (let i = 0; i < plinkos.length; i++) {
            plinkos[i].show();
        }

        for (let i = 0; i < bounds.length; i++) {
            bounds[i].show();
        }
    }

    p.draw = function () {
        Engine.update(engine);

        frame++;

        // Every 24 frames
        if (frame % 24 == 0) {

        }

        // Draw black balls in ball positions
        for (let i = 0; i < particles.length; i++) {
            //particles[i].show(ballPositions[i]);
        }

        let nextBallPositions = {}

        for (let i = 0; i < particles.length; i++) {
            nextBallPositions[i] = particles[i].show();
            particles[i].isOffScreen()

            // Add this code to emit the ball's position
            if (p.movies && p.movies.length) {
                const segmentIndex = p.getSegmentIndex(particles[i]);
                if (segmentIndex >= 0 && segmentIndex < p.movies.length) {
                    // Dispatch custom event with the movie index
                    const event = new CustomEvent('ballPositionUpdate', {
                        detail: { movieIndex: segmentIndex }
                    });
                    window.dispatchEvent(event);
                }
            }

            /*
if (p.movies.length && this.checkIfStopped(particles[i])) {
const segmentIndex = this.getSegmentIndex(particles[i]);
console.log(segmentIndex)
console.log(`Particle stopped in segment: ${segmentIndex}`);
console.log(p.movies[segmentIndex]);
setTimeout(() => {
    if (!alerted) {
        alert(`You will watch ${p.movies[segmentIndex].title}`);
    }
    alerted = true;

   // Remove the particle from the world
    particles.splice(i, 1);
    i--;
}, 1000);
}
            */
        }

        ballPositions = nextBallPositions

    };

    function Particle(x, y, rad) {
        this.r = 255;
        this.g = 255;
        this.b = 255;

        const options = {
            isStatic: false,
            mass: 10,
            density: 1,
            restitution: 1,
            friction: 0
        };

        this.body = Bodies.circle(x, y, rad, options);
        this.body.label = "particle";
        this.r = rad;
        this.color = [255, 255, 255];
        this.type = Math.floor(Math.random() * 3);
        this.body.targetAngle = this.body.angle;

        World.add(world, this.body);
    }

    Particle.prototype.isOffScreen = function () {
        const { x, y } = this.body.position;

        if (x < 0) {
            Matter.Body.setPosition(this.body, { x: p.width, y });
        } else if (x > p.width) {
            Matter.Body.setPosition(this.body, { x: 0, y });
        }

        if (y > p.height) {
            return true;
        }

        return false;
    };

    Particle.prototype.updateRotation = function () {
        const currentAngle = this.body.angle;
        const targetAngle = this.body.targetAngle;
        const angleDifference = targetAngle - currentAngle;

        // Apply a small fraction of the difference to the current angle
        this.body.angle += angleDifference * 0.1;
    };


    Particle.prototype.show = function (oldPosition = false) {
        const pos = oldPosition ? oldPosition : this.body.position;


        p.noStroke();

        if (oldPosition) {
            p.fill(21);
        }
        else {
            p.fill(255, 255, 255);
        }

        p.push();
        p.translate(pos.x, pos.y);
        p.ellipse(0, 0, this.r * 2 + (oldPosition ? 0 : 0))
        p.pop();

        /*
        // Draw a new circle on top of the particle
        p.fill(21);
        p.push();
        p.translate(pos.x - 6, pos.y);
        p.ellipse(0, 0, 3);
        p.pop();

        p.push();
        p.translate(pos.x + 6, pos.y);
        p.ellipse(0, 0, 3);
        p.pop();

        p.push();
        p.translate(pos.x, pos.y + 1);
        p.rect(-2.5, 0, 5, 1);
        p.pop();
        */

        return { x: pos.x, y: pos.y }
    };

    // ======================================================
    //           Plinko.js
    // ======================================================

    function Plinko(x, y, r) {
        const options = {
            isStatic: true,
            density: 1,
            restitution: 1,
            friction: 0,
        };
        let color = palette[Math.floor(Math.random() * palette.length)];
        this.color = [color._rgb[0], color._rgb[1], color._rgb[2]];
        this.body = Bodies.circle(x + (Math.random() * 12) - 6, y + (Math.random() * 12) - 6, r, options);
        this.body.label = "plinko";
        this.body.frequency = (Math.random() * 300) + (y * 0.5);
        this.r = r;
        World.add(world, this.body);
    }

    Plinko.prototype.show = function () {
        p.fill(this.color);
        // stroke(255);
        p.noStroke();
        const { x, y } = this.body.position;
        p.push();
        p.translate(x, y);
        p.ellipse(0, 0, this.r * 2);
        p.pop();
    };

    // ======================================================
    //           Boundary.js
    // ======================================================

    function Boundary(x, y, w, h) {
        const options = {
            density: 1,
            friction: 1,
            isStatic: true
        };
        this.body = Bodies.rectangle(x, y, w, h, options);
        this.body.label = "boundary";
        this.w = w;
        this.h = h;
        World.add(world, this.body);
    }

    Boundary.prototype.show = function () {
        p.fill('#222');
        // stroke(255);
        p.noStroke();
        const { x, y } = this.body.position;
        p.push();
        p.translate(x, y);
        p.rectMode(p.CENTER);
        p.rect(0, 0, this.w, this.h);
        p.pop();
    };
};

export default {
    name: "Plinko",
    data() {
        return {
            widthPercentages: null,
            myp5: null,
            activeMovieIndex: null
        };
    },
    props: {
        movies: {
            type: Array,
            required: true,
        },
        list: {
            type: String,
            required: true
        },
        advancedOptions: {
            type: Object,
            required: true
        }
    },
    computed: {
        thisList() {
            const parts = this.list.split('/');
            return 'https://letterboxd.com/' + parts[0] + '/list/' + parts[1];
        }
    },
    methods: {
        dropBall() {
            this.myp5.newParticle();
        },
        updateActiveMovie(index) {
            this.activeMovieIndex = index;
        }
    },
    mounted() {
        this.myp5 = new p5(sketch, this.$refs.myCanvas);
        // Wait for movies to have length
        this.myp5.populate(this.movies);
        setTimeout(() => {
            this.widthPercentages = this.myp5.getWidthPercentages()
        }, 2000)

        // Add event listener for ball position updates
        window.addEventListener('ballPositionUpdate', (e) => {
            this.updateActiveMovie(e.detail.movieIndex);
        });
    },
    beforeUnmount() {
        // Clean up event listener
        window.removeEventListener('ballPositionUpdate', this.updateActiveMovie);
    }
};
</script>

<style>
canvas {
    width: 100vw;
    height: auto;
    margin-block-start: 3em;
    border: 3px solid #222;
    border-top: 0;
    position: relative;
    z-index: 2;
}

.images {
    position: absolute;
    top: calc(100% - 10px);
    left: -1px;
    width: 100%;
    height: auto;
    display: flex;
    padding-bottom: 15em;
    z-index: 1;
}

.parent {
    position: relative;
    width: fit-content;
    margin: auto;
}

.image a {
    display: block;
}

.image img {
    width: calc(100% - 5px);
    height: auto;
    margin-inline: 2.5px;
    object-fit: cover;
    filter: brightness(0.5);
    transition: filter 150ms linear;
}

.image.active img {
    filter: brightness(1.5);
}

.drop {
    margin-block-start: 2em;
    appearance: none;
    -webkit-appearance: none;
    width: 100%;
    background: transparent;
    border: 2px solid #222;
    color: #fff;
    display: inline-flex;
    max-width: 500px;
    padding: 1rem;
    font-weight: bold;
    text-transform: uppercase;
    text-align: center;
    font-size: 1.5em;
    justify-content: center;
}

.drop:active {
    transform: scale(0.95);
}

.movie-name {
    white-space: nowrap;
    text-align: right;
    display: block;
    transform-origin: 0% center;
    width: 100%;
    transform: translateX(50%) rotate(90deg);
    opacity: 0.4;
    transition: all 150ms linear;
    display: flex;
    gap: 1em;
}

.movie-name a {
    opacity: 0.5;
    font-weight: bold;
    text-decoration: none;
    color: currentColor;
    display: inline-block;
    position: relative;
}

.movie-name a::before {
    content: '';
    display: inline-block;
    position: absolute;
    width: 120%;
    height: 150%;
    top: -25%;
    left: -10%;
}

.movie-name a.clicked {
    opacity: 1;
    color: oklch(75.23% 0.209 144.64);
}

.dark .image.active .movie-name {
    color: yellow;
}

.image.active .movie-name {
    opacity: 1;
    font-weight: bold;
    text-shadow: 0px 0px 10px currentColor;
}
</style>
