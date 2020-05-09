## The What

This example consists of a **HTTP server** who accepts request to generate PDF from given URLs.

It works as one would expect, except for the fact that **most of the time** it hangs forever after a few requests. But **not always**. Sometimes it works for quite a few requests in a row.

On the flip side, when using "load.windowStatus" property, **most of the time** it hangs forever. Again, **not always**. Sometimes it works for quite a few requests.

## The Why

OK. But if **it is so buggy**, why is it here? Because it serves as a lab to try and find a solution for the given problem. So people can work from the same ground and build up a solution.

You might help us find the solution. Or it might be that you already know and can share the solution -- if there is one.

## The Maybe

It might be the fact that **QT does have to work on the main thread**, but I've tried different implementations with no success.

Spicing up the nature of the problem, there is a project [cheesyd](https://github.com/leandrosilva/cheesyd) which does something somewhat similar in C++ and it works a charm with zero hangs.