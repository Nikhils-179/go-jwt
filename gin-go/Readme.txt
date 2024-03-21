	• High Performance : Built with speed in mind. Ideal for handling high-traffic backend systems .
	• Minimalistic and Light weight : Design philosophy to keep things simple and minimal.
	• Fast Router : It is highly optimized and can quickly handle routing tasks
	• Middleware support : provides robust middleware system,allwoing developers to extend functionalit such as authentication,logging 
	
	Installation :
	1. Create a New Go module 
		- Mkdir gin-be
		- Cd gin-be
		- Git mod init github.com/your-username/gin-be
	2. Installing Gin Package
		- Go get -u github.com/gin-gonic/gin
	 
	

Middleware in Gin
	• Middleware in Gin is an essential component that intercepts HTTP requests and responses 
	• They perform pre-processing before a request reaches the designated route-handler or post-processing tasks before the response is sent to client
	• 

	• In this we defined a loggerMiddleware function that caculates the duration of each request and logs thee method.
	• Router.Use() to apply our custom logger to all routes 

Creating custom middleware 

	• Developer often need to implement custom middleware for project specific requirements.
	• Custom middle ware can handle task like authentication , data validation , rate limiting 
	• 


Routing and Grouping
	
	• Routing is mapping incoming HTTP requests to specific route handlers
	• The router matches the URL path and HTTP method of the request to find appropriate handler to exeute
	• 

	• Route Groups allows you to group related routes , which makes the code more organized and easier to maintain.
	• Line no . 54 , 46,66
	• 


Controllers and Handlers

As Backend grows , handling all business logic in route handlers becomes unweirdly . To improve code organization and maintainability , Gin encourages using controllers to handle business logic seperatly from route handlers 

SEPERATING BUISNESS LOGIC FROM CONTROLLERS


 Create handler/controller  ( define end point)
 Create API for endpoint


