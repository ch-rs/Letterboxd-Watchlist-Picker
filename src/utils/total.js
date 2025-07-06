export function total(arr) {
    return arr.reduce((partialSum, a) => partialSum + a, 0)
}