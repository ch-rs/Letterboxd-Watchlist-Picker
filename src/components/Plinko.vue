<template>
    <div @click="dropBall">
        <div ref="myCanvas"></div>
    </div>
</template>

<script>
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
        cols = 11,
        spacing,
        rows = 9,
        particleSize = 10,
        movies = [],
        slotWidth,
        imageObjects = [],
        stopped = false,
        plinkoSize = 14;

    p.setup = function () {
        const c = p.createCanvas(612, 700);
        c.parent(parent);

        engine = Engine.create();
        world = engine.world;
        world.gravity.y = 2;

        spacing = p.width / cols;

        Matter.Events.on(engine, "afterUpdate", function () {
            for (let i = 0; i < particles.length; i++) {
                // Check if the particle's speed is close to zero
                if (stopped == false && particles[i].body.speed < 0.01) {
                    stopped = true;
                }
            }
        });
    };

    p.populate = function (newMovies) {
        movies = newMovies;
        for (let i = 0; i < movies.length; i++) {
            const x = slotWidth * (i + 1) - slotWidth;

            // Create a new ImageObject at the calculated position
            const obj = new ImageObject(
                x,
                p.height - 100,
                movies[i].image_url,
                slotWidth,
            );

            imageObjects.push(obj);
        }

        // Calculate the width of each slot
        slotWidth = p.width / movies.length;

        p.createBoundaries();
        p.createPlinkos();
        p.populate();
    };

    function ImageObject(x, y, imageUrl, slotWidth) {
        this.x = x;
        this.y = y;
        this.image = p.loadImage(imageUrl);

        this.show = function () {
            // Specify the width and height of the image
            p.image(
                this.image,
                this.x,
                this.y,
                slotWidth,
                (slotWidth * this.image.height) / this.image.width,
            );
        };
    }

    p.newParticle = function () {
        const p = new Particle(random(100, 600), 0, particleSize);
        particles.push(p);
    };

    p.createBoundaries = function () {
        let b = new Boundary(p.width / 2, p.height + 50, p.width, 100);
        bounds.push(b);

        for (let i = 0; i < cols + 1; i++) {
            const x = i * slotWidth;
            const h = 100;
            const w = 10;
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
                const p = new Plinko(x, y, plinkoSize);
                plinkos.push(p);
            }
        }
    };

    p.draw = function () {
        p.background(51);
        Engine.update(engine);

        for (let i = 0; i < imageObjects.length; i++) {
            imageObjects[i].show();
        }

        for (let i = 0; i < particles.length; i++) {
            particles[i].show();
            if (particles[i].isOffScreen()) {
                World.remove(world, particles[i].body);
                particles.splice(i, 1);
                i--;
            }
        }

        for (let i = 0; i < plinkos.length; i++) {
            plinkos[i].show();
        }

        for (let i = 0; i < bounds.length; i++) {
            bounds[i].show();
        }
    };

    function Particle(x, y, r) {
        this.r = 255;
        this.g = 255;
        this.b = 255;
        const options = {
            isStatic: false,
            mass: 0,
            density: 1,
            restitution: 0.5,
            friction: 1,
        };
        x += random(-1, 1);
        this.body = Bodies.circle(x, y, r, options);
        this.r = r;
        World.add(world, this.body);
    }

    Particle.prototype.isOffScreen = function () {
        const { x, y } = this.body.position;
        return x < -50 || x > p.width + 50 || y > p.height + 50;
    };

    Particle.prototype.show = function () {
        p.noStroke();
        p.fill(this.r, this.g, this.b);
        const pos = this.body.position;
        p.push();
        p.translate(pos.x, pos.y);
        p.ellipse(0, 0, this.r * 2);
        p.pop();
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
        this.color = random(80, 175);
        this.body = Bodies.circle(x, y, r, options);
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
            isStatic: true,
        };
        this.body = Bodies.rectangle(x, y, w, h, options);
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
    },
    methods: {
        dropBall() {
            this.myp5.newParticle();
        },
    },
    mounted() {
        this.myp5 = new p5(sketch, this.$refs.myCanvas);
        this.myp5.populate(this.props.movies);
    },
};
</script>

<style>
canvas {
    width: 100vw;
    height: auto;
    margin-block-start: 3em;
}
</style>
