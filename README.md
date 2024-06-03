# Planny

Planny is a backend application developed in Go (Golang) that offers a simple features for students to effectively
manage their daily plans. This project was created as part of the VatanSoft internship case, aiming to showcase the
candidate's skills in backend development using Golang, GORM, MySQL, and the Echo web framework.

## Technologies

- Golang
- GORM (Go Object-Relational Mapping)
- MySQL
- Echo (Web Framework)

## Implemented Features

- **Record Student Daily Plans**: Students can register their daily plans.
- **Plan Scheduling**: Students can record a plan with specific days and time ranges.
- **Plan States**: Plans have states such as canceled, done, and in progress.
- **Plan Updates and Deletion**: Students can update and deleted the plans.
- **Conflict Checking**: Checking if there is another plan during the same date and time range when adding a new plan.
- **Student Registration and Information Update (Optional)**: Allow students to register and update their information.

## API Endpoints

All endpoints are prefixed with `/api/v1`. For more
info [Postman Workspace](https://www.postman.com/planetary-moon-654796/workspace/planny/collection/32427111-a2852ce0-76f1-46bf-93f0-0465dbded2f7?action=share&creator=32427111).

### Authentication

| Method | Path          | Description               |
|--------|---------------|---------------------------|
| POST   | /register     | Register a new student    |
| POST   | /login        | Login an existing student |
| POST   | /renew_access | Renew Access Token        |

### Plans

| Method | Path       | Description                      |
|--------|------------|----------------------------------|
| GET    | /plans     | Get all plans belongs to student |
| POST   | /plans     | Create a new plan                |
| PATCH  | /plans/:id | Update a plan                    |
| DELETE | /plans/:id | Delete a plan                    |

## License

This project is licensed under the [Apache License](./LICENSE).

## Authors

- [Samet Demir](https://github.com/asdsec)
