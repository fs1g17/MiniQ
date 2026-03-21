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

### Reactivity

Currently I'm making the worker poll the queue
that's not a terribly good idea - how can I improve it?

I was thinking about this approach:

- queue gets a new job pushed
- the QUEUE then checks whether any workers are currently working:
  - if any worker is available, then we call `go worker.Perform(job)`
  - if no workers are available, we do nothing
  - then, when a worker finishes, check the queue once from INSIDE the worker
  - this way, the workers and the queue can sleep and only react
  - there could be a race condition if 2 workers check a queue while something is being enqueued,
    which is currently empty, they sleep, but queue thinks they're both busy, can avoid by using mutex smartly

what's a good way of tying that stuff in together?
i want a "knock-on" effect to be triggered when something is enqueued in the queue - are there events in go?
welp apparently folks are using messaging queues for events - LOL
i guess there's nothing wrong in making the queue do the checking for now
that does mean the queue is kinda coupled to the workers
i'll do it this way and think about a better approach in the meantime

so what if i keep the queue as the main controller

- when a worker frees up, we notify the queue that a worker freed up
- if there's jobs in the queue, the queue then invokes the worker
- if theres no jobs, everything runs as is until all workers free up
- as soon as a new job comes in, the queue assigns it to the first worker

### Concurrently setting "busy"

So I've noticed something.

If I have a struct like the worker:

type Worker struct {
Status WorkerStatus
mu sync.Mutex
}

even if i have some getters and setters that rely on the mutex:

func (w \*Worker) SetStatus(ws WorkerStatus) {
w.mu.Lock()
defer w.mu.Unlock()
w.Status = ws
}

func (w \*Worker) GetStatus() WorkerStatus {
w.mu.Lock()
defer w.mu.Unlock()
return w.Status
}

I have to set the status in the thread that calls the goroutine
if I have a perform function:

func (w \*Worker) Perform() {
w.SetStatus(Busy)
}

this doesn't make it concurrent - I had this example, I had 2 calls to miniQ.addJob - which looped over the Worker slice to find workers in the Idle state. But both times, it found worker0 - which is wrong. Setting the status in the calling thread was what fixed it. No idea why, gonna have to read about how goroutines are actually executed.
