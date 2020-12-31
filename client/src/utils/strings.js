export function snakeToSentenceCase(snake) {
    return snake
        .split("_")
        .map(word => word[0].toUpperCase() + word.slice(1))
        .join(" ")
}