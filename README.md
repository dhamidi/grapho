# Description

γράφω is Greek for "I write".  This piece of software is going to power my yet to be blog.  There are many solutions to the problem "I want to publish articles on the web", but this one is mine.  Besides solving this problem for me, γράφω serves as a test subject for the following:

- deploying and monitoring a side-project in production: haven't done this yet and it is about time.
- [Domain Events][domain-events] / [Event Sourcing][event-sourcing]: promise me the following:
  - rollback of changes
  - good performance
  - loosely coupled modules

  If this works in practice, it changes a lot of accumulated "wisdom".
  It might be just a better way to develop software compared to what
  I've been doing so far.
- [No relational database in the center of the application][nodb]

# Scope / Requirements

- γράφω is an application, not a platform.  As such, it is single-user only.
- Authentication / authorization is left to the protocol layer.
- User-generated content (e.g. comments) requires approval before being published.

# Use cases

- [ ] M1 Draft post
- [ ] M1 List drafts
- [ ] M1 Publish post 
- [ ] M1 List all posts
- [ ] M1 Show post
- [ ] M2 Comment on post
- [ ] M2 Approve comment
- [ ] M2 Amend comment

[domain-events]: http://martinfowler.com/eaaDev/DomainEvent.html
[event-sourcing]: http://martinfowler.com/eaaDev/EventSourcing.html
[nodb]: http://blog.8thlight.com/uncle-bob/2012/05/15/NODB.html

# LICENSE

See [LICENSE](./LICENSE).
