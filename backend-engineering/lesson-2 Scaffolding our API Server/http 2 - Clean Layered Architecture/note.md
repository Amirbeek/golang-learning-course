## Dependency Inversion Principle (DIP)
you are injecting the dependencies in your layers. you font directly call them. 
Why? it promotes loose coupling and makes it easier to test our program.


## Adapting to change:
It means being able to handle new situations, challenges, or environments smoothly.
In work or life, things often don’t go as planned — new tools, rules, or problems can appear.
When you adapt, you stay calm, learn quickly, and find new ways to succeed. For example,
if your company switches to a new software or your project goals change,
you adjust your methods instead of giving up. This skill shows flexibility,
problem-solving, and a positive attitude toward growth.


## Focus on Business Value
It means working on things that bring real benefit to the business or customers. 
Instead of just doing tasks, you think about how your work helps the company grow,
saves time, earns money, or improves user satisfaction.
Every decision and effort should support the main business goals and deliver value to the customer.


1) Transport  -- http handlers
2) Service <- Repository 
3) Storage

## Injected Dependencies:
Bu shuni anglatadi: dastur implementations (aniq sinflar yoki structlar) ga emas, balki abstractions (ya’ni interface) ga tayanadi.
Bunda biz interface orqali bog‘lanamiz, bu esa kodni flexible, reusable, va easy to test qiladi.
Masalan, agar kelajakda biror qismni o‘zgartirmoqchi bo‘lsak, interface tufayli boshqa joylarga ta’sir qilmaymiz.


# Service Layer
## Service 
app.CreateAndInviteUsers()

  |
  |
  v
## Repository
* CreateUser()
* CreateInvite()

  |
  |
  |
  v
DB