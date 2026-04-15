# Project Requirements Specification

---

## General Logic and Backend Updates

### 1. Data Schema Enhancements
* **Project Entity**: Add a date_approved field.
* **Student-Project Status**: Implement three distinct states:
    * PENDING_APPROVAL (Waiting)
    * SCHOOL_APPROVED (School has verified)
    * LOCAL_APPROVED (Local authority has verified)
* **User Access Control**:
    * Add a status to the User account (Approved vs. Unapproved).
    * **Restriction**: Unapproved users are restricted to viewing their own profile only. They cannot view system data or perform any edit operations.

### 2. System Mechanics
* **Email System**: Automated notifications for student application outcomes (Accepted or Rejected).
* **Anti-Spam**: Implement a mechanism to prevent students from spamming project registrations. (not implement for now)
* **Architecture Justification**: (when read to implement, you can ignore this)
    * **Question**: Why store Student, School, and Local data in a single table?
    * **Answer**: Centralization. Since all roles share core authentication data (email, password, basic profiles), a unified table simplifies the login process, session management, and Role-Based Access Control (RBAC).

### 3. Enumerations (Enums)
* Add this to StudentProject, the flow will be applied to project where `data_approved` is not nil. The meaning of this table is the application of a student to join a project from a community:
  * SCHOOL_PENDING: when student registers to a project, wait for school approval
  * SCHOOL_REJECT: the student's application is rejected by school
  * COMMUNITY_PENDING: school accepts -> waiting for the approval of the community
  * COMMUNITY_REJECT the student's application is rejected by community
  * APPROVED: the student's application is approved

---

## User Interface (UI) Requirements
Below is the UI, you can use it as a referene for the backend implementation

### General UI Components
* **Account Status**: Display a badge or status indicating if the user is approved.
* **Settings Page**: A dedicated section for users to update Contact Information and Avatar only.

### Role-Specific Views
* **School**:
    * View all projects.
    * View participant lists for each project to approve students.
    * Access a review page for projects sent to the school.
* **Local (Community)**:
    * **Home**: List of projects created by the local user.
    * **Create Project**: Form to submit new initiatives.
    * **Project Details**: View list of participating students (Note: Local role can only Reject students).
* **Student**:
    * View a list of projects available to their specific school.
    * Action button to Register for projects.
* **Admin**:
    * Master view of all projects.
    * Management page for pending and active users.
    * Affiliation management page (with popup/modal for adding/removing new ones).

---

## Implementation notice 
* only student can apply
* Make the approval endpoint to be flexible, that is, accepts multiple project
* Add 2 endpoint that serves static file for the profile picture and project banner
* Add description field for affiliation
* Make `is_active = false` user can't login 
* Be explicit about which field is optional or mandatory.
* Make sure each user can only view its appropriate content.
* Make sure to check the form registration date of the application of a student.
* create an endpoint for each role to view the applications

## Testing
* Create unit test and e2e test. The e2e test will be like this.
  * Admin account will be pre-created.
  * Admin add 2 affilations (1 school, 1 community). Total there are 3 affiliations (including the admin affiliation)
  * User, school, community create account.
  * Admin approves the accounts.
  * community creates a project
  * community view the project -> see it hasn't been approved
  * student view its affiliation's projects -> see nothing
  * school approves the project
  * community view the project -> see it has been approved
  * student view its affilation's projects -> see the community newly created project.
  * school opens the project -> view its registration list and see no applicant
  * student register to the project in the correct form date.
  * school opens the project -> view its registration list and see an application from a student. Can also see the status (our defined ENUM above)
  * community opens the project -> see no applicant
  * school approves the student applicant.
  * student views their application list -> see it status is COMMUNITY_PENDING
  * commnuity opens the project -> see an applicant
  * community accepts the application
  * student views their application list -> see its status is APPROVED
  * school virews their application list -> see its status is APPROVED
