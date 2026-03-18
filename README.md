# MiniQ

This is my little Queue project. The purpose of MiniQ is to learn about Go.

## Thought pattern

How would I implement a worker?
A worker is something that executes a piece of code.
I guess the actual "work" could be an anonymous function.

In JavaScript, I would've done something like this:
class Worker {
constructor(work) {
this.work = work;
}
}
Since there's no typing, the function "work" would accept any type, and it would just break at runtime if something is wrong.

in TypeScript, which is closer to Go, I guess I'd use a generic type to define the arguments to the work function, this way a worker can only ever work on an appropriate queue. I guess we could then also type the Queue.

class Worker<T> {
constructor(work: (args: T) => void) {
this.work = work;
}
}

I'm not sure about the syntax 100% but maybe something like this. I'll give this a try now.
