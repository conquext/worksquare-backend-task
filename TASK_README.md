
# Worksquare Senior Backend Vetting Task

Welcome to the vetting assessment for Senior Backend Developers at Worksquare.

This task evaluates your ability to design and implement scalable RESTful APIs using clean architecture and modern engineering practices.

---

## Objective

Build a RESTful API to serve housing listings from a `listings.json` file. Your goal is to demonstrate API architecture, security, pagination, and clean code organization.

---

## Supported Stacks

Choose one of the following:

* Node.js (Express)
* Laravel (PHP)
* Django or Flask (Python)

---

## Core Requirements

Your API should:

* Load listing data from a `listings.json` file (shared separately)
* Return a **paginated** list of listings
* Allow **filtering** by `location` and `property type`
* Allow **retrieving** a listing by `ID`
* Implement **JWT authentication** (protect at least one route)
* Use **rate limiting** (e.g., 100 requests per hour per IP)
* Log **incoming requests** via middleware
* Handle **errors** with standardized responses
* Support pagination with `page` and `limit` query parameters
* Include **Swagger/OpenAPI documentation**
*  **Dockerize** the setup
* Add **unit or integration tests**

---

## Architecture & Approach

Please include the following (in a `docs/` folder or in the README):

* Overview of your architecture and folder structure
* Database/data model explanation (if applicable)
* Your authentication & security approach
* API design strategy (REST principles, middleware, error flow)
* Trade-offs and decision-making process

---

## Documentation

Ensure your `README.md` includes:

* Setup instructions (environment, dependencies, run commands)
* Tools and libraries used
* Swagger documentation link or UI preview
* "Code Notes" explaining your development approach
* Screenshots or API usage examples (optional)

---

## What We’re Looking For

* Clean and modular code
* Proper use of REST conventions
* Scalable architecture
* Solid Git practices and commit messages
* Good developer experience (DX)

---

## Submission Guidelines

1. Fork this repository
2. Scaffold your backend project inside
3. Commit regularly with meaningful messages
4. Push to your forked repo
5. Email us:

   * GitHub repo link
   * Notes (if any)

---

## Deadline

Submit your completed task within **48 hours** of receiving it.

---

## Listings Data

You’ll receive the `listings.json` file separately. Please create this file in your local project and use it as your API data source.
Note! you can use Ai for your task but ensure to document where it was used
