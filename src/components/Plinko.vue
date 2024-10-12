<template>
    <div>
        <button class="drop" @click="dropBall">Drop!</button>
        <div class="parent" ref="myCanvas">
            <div v-if="movies" class="images">
                <div v-for="movie in movies" class="image">
                    <img :src="movie.image_url" />
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import ding from '../utils/ding';

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
        rows = 15,
        particleSize = 12,
        slotWidth,
        ballPositions = {},
        alerted = false,
        plinkoSize = 14;

    p.setup = function () {
        const c = p.createCanvas(612, 1100);
        c.parent(parent);

        engine = Engine.create();
        world = engine.world;
        world.gravity.y = 2;

        spacing = p.width / cols;

        Matter.Events.on(engine, "collisionStart", function (event) {
            for (let i = 0; i < event.pairs.length; i++) {
                let pair = event.pairs[i];
                let velocityA = pair.bodyA.velocity;
                let velocityB = pair.bodyB.velocity;

                console.log(`Velocity of bodyA: x=${velocityA.x}, y=${velocityA.y}`);
                console.log(`Velocity of bodyB: x=${velocityB.x}, y=${velocityB.y}`);

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

        p.background(21);

    };

    p.populate = function (newMovies) {
        cols = newMovies.length;
        p.movies = newMovies;

        // Calculate the width of each slot
        setTimeout(() => {

            slotWidth = p.width / newMovies.length;

            // Loop through p.movies give each an X position based on priority value
            for (let i = 0; i < p.movies.length; i++) {
                p.movies[i].x = i * slotWidth;
            }

            p.createPlinkos();
            p.createBoundaries();

            p.preDraw()
        }, 1000);
    };

    p.newParticle = function () {
        alerted = false
        const part = new Particle(random(5, p.width - 5), 0, particleSize);
        particles.push(part);
    };

    p.createBoundaries = function () {
        let b = new Boundary(p.width / 2, p.height + 50, p.width, 100);
        bounds.push(b);

        for (let i = 0; i < cols + 1; i++) {
            const x = i * slotWidth;
            const h = 60;
            const w = 5;
            const y = p.height - h / 2;
            b = new Boundary(x, y, w, h);
            bounds.push(b);
        }
    };

    p.createPlinkos = function () {
        for (let j = 0; j < rows; j++) {
            for (let i = 0; i < cols + 1; i++) {
                let x = i * spacing;
                if (j % 2 == 0) {
                    x += spacing / 2;
                }
                const y = spacing + j * spacing;
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
        const x = particle.body.position.x;
        return Math.floor(x / sw);
    }

    p.preDraw = function () {
        for (let i = 0; i < plinkos.length; i++) {
            plinkos[i].show();
        }

        for (let i = 0; i < bounds.length; i++) {
            bounds[i].show();
        }
    }

    p.draw = function () {
        Engine.update(engine);

        // Draw black balls in ball positions
        for (let i = 0; i < particles.length; i++) {
            particles[i].show(ballPositions[i]);
        }

        let nextBallPositions = {}

        for (let i = 0; i < particles.length; i++) {
            nextBallPositions[i] = particles[i].show();
            particles[i].isOffScreen()

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
            mass: 0,
            density: 1,
            restitution: 1,
            friction: 1,
        };

        this.body = Bodies.circle(x, y, rad, options);
        this.body.label = "particle";
        this.r = rad;
        this.color = [255, 255, 255];

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

    Particle.prototype.show = function (oldPosition = false) {

        p.noStroke();

        if (oldPosition) {
            p.fill(21);
            p.push();
        }
        else {
            p.fill(255, 255, 255);
            p.push();
        }

        p.push();

        const pos = oldPosition ? oldPosition : this.body.position;
        p.translate(pos.x, pos.y);

        p.ellipse(0, 0, this.r * 2 + (oldPosition ? 1.3 : 0))
        p.pop();

        return {x: pos.x, y: pos.y}
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
        this.color = [random(80, 150), random(80, 150), random(80, 150)];
        this.body = Bodies.circle(x, y, r, options);
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
        p.fill(128);
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
            myp5: null,
        };
    },
    props: {
        movies: {
            type: Array,
            required: true,
        },
        advancedOptions: {
            type: Object,
            required: true
        }
    },
    methods: {
        dropBall() {
            this.myp5.newParticle();
        },
    },
    mounted() {
        this.myp5 = new p5(sketch, this.$refs.myCanvas);
        // Wait for movies to have length
        this.myp5.populate(this.movies);
    },
};
</script>

<style>
canvas {
    width: 100vw;
    height: auto;
    margin-block-start: 3em;
}

.images {
    position: absolute;
    top: calc(100% - 10px);
    left: -1px;
    width: 100%;
    height: auto;
    display: flex;
    background: #808080;
}

.parent {
    position: relative;
    width: fit-content;
    margin: auto;
}

.image img {
    width: calc(100% - 5px);
    height: auto;
    margin-inline: 2.5px;
    object-fit: cover
}

.drop {
    margin-block-start: 2em;
    appearance: none;
    -webkit-appearance: none;
    width: 100%;
    background: #000;
    color: #fff;
    display: inline-flex;
    border: 0;
    max-width: 500px;
    padding: 1rem;
    font-weight: bold;
    text-transform: uppercase;
    text-align: center;
    font-size: 1.5em;
    justify-content: center;
}
</style>
