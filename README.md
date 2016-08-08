# gits
#### a simple git deployment

This project was made out of the need for simplicity. Moving away from a very
clunky Jenkins install to something much lighter. All it does is listen for
GitHub hooks, and updates the repository.

As of right now it's not _very_ customizable, but feel free to add configuration
options as you see fit. Just try to keep it from turning into a Jenkins install
`;)`.

Right now you configure GitHub hooks manually (with a secret) and gits will
listen to `create` and `push` events. It will then create the directory, clone
the repository, checkout the branch, and optionally run a script.

### Installation
gits does not require any compiling, just make sure you have a newer version of
`node` and run `npm start`. Just make sure to configure it first `:P`.

### License
MIT
