On Open/Close
=============

The "network" layer `OnOpen` and `OnClose` attributes function exactly the same as the "generic" 
layer functions, with the exception that the driver value passed into the user provided 
functions is of type `network.Driver` rather than `generic.Driver`. Typically, users will only 
use one of the two "flavors" of OnX functions, but if you wanted to you could of course use both!