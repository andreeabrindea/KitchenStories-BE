using ErrorOr;
using GourmetStories.Contracts.GourmetStories;
using GourmetStories.Services;
using GourmetStories.Models;
using Microsoft.AspNetCore.Mvc;

namespace GourmetStories.Controllers;

public class UsersController : ApiController
{
    private readonly IUserService _userService;

    public UsersController(IUserService userService)
    {
        _userService = userService;
    }
    [HttpPost]
    public IActionResult CreateRecipe(CreateUserRequest request)
    {
        var user = Models.User.Create(
            request.Name,
            request.Username,
            request.Password,
            request.Email,
            Guid.NewGuid()
        );

        if (user.IsError)
        {
            return Problem(user.Errors);
        }
        var createRecipeResult = _userService.CreateUser(user.Value);
        return createRecipeResult.Match(
            _ => CreatedNewUser(user.Value),
            Problem);
    }
    
    [HttpGet("{id:guid}")]
    public IActionResult GetUser(Guid id)
    {
        ErrorOr<User> getUserResult = _userService.GetUser(id);
        return getUserResult.Match(
            user => Ok(MapUserResponse(user)),
            Problem);
    }

    [HttpPut("{id:guid}")]
    public IActionResult UpsertUser(Guid id, UpsertUserRequest request)
    {
        var user = Models.User.Create(
            request.Name,
            request.Username,
            request.Password,
            request.Email,
            id
        );

        if (user.IsError)
        {
            return Problem(user.Errors);
        }
        
        ErrorOr<UpsertUserResult> updateRecipeResult = _userService.UpsertUser(user.Value);
        return updateRecipeResult.Match(
            updated => NoContent(),
            Problem);
    }

    [HttpDelete("{id:guid}")]
    public IActionResult DeleteRecipe(Guid id)
    {
        ErrorOr<Deleted> deleteUserResult = _userService.DeleteUser(id);
        return deleteUserResult.Match(
            _ => NoContent(),
            Problem);
    }

    private static User MapUserResponse(User user)
    {
        return Models.User.Create(
            user.Name,
            user.Username,
            user.Password,
            user.Email,
            user.Id
        ).Value;
    }

    private CreatedAtActionResult CreatedNewUser(User user)
    {
        return CreatedAtAction(
            actionName: nameof(GetUser),
            routeValues: new { id = user.Id },
            value: MapUserResponse(user)
        );
    }
}