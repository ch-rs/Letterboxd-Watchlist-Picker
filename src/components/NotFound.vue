<template>
    <div>
        <h2>Nothing Found</h2>
        <p v-if="status == 'no-intersect'">
            The intersection between those two lists is empty.
        </p>
        <p v-else-if="status == 'possibly-ignored'">
            There were no films found in that list. It may be empty, private, or
            only contain films ignored by settings (check Advanced Options
            above).
        </p>
        <p v-else-if="status == 'timeout'">
            Sorry your list was too powerful and we timed out. Try another list.
        </p>
        <p v-else>
            Sorry, that watchlist is empty, private, or doesn't exist at all.
        </p>
    </div>
</template>

<script>
export default {
    name: "NotFound",
    props: ["status"],
};
</script>

<style scoped>
h2 {
    margin-top: 5rem;
    font-size: 1.2rem;
}

p {
    max-width: 60ch;
    margin: 0rem auto 1rem;
    padding: 0 2rem;
}

#not-found-gif {
    margin: 30px 0;
    border-radius: 4px;
}
</style>
