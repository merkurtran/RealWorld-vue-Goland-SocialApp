describe("Auth Component Tests", () => {
    beforeEach(() => {
        cy.visit('/Auth')
    })

    describe("Page Layout", () => {
    it('should display sigin and signup cards', () => {
        // check for bot cards exists
        cy.get('.q-card').should('have.length', 2)

        // check titles
        cy.contains('Signin').should('be.visible')
        cy.contains('Signup | Create New Account').should('be.visible')
    })

    it('should have proper layout structure', () => {
        cy.get('.col-5').should('exist')
        cy.get('.col-7').should('exist')
        cy.get('.row').should('exist')
    })
})

    describe('Sigin Form', () => {
        it('should display all signin form elements', () => {
            cy.get('.col-5 input').should('have.length', 2)

            // check for button
            cy.contains('sign in').should('be.visible')
            cy.get('.col-5 button[type="submit"]').should('be.visible')
        })

        it('should allow typing in signin inputs', () => {
            cy.get('.col-5 input').eq(0)
                .type('test@example.com').should('have.value', 'test@example.com')

            cy.get('.col-5 input').eq(1)
                .type('password123').should('have.value', 'password123')
        })

        it('should have password input type', () => {
            cy.get('.col-5 input').eq(1)
                .should('have.attr', 'type', 'password')
        })
    })
    describe('Signup Form', () => {
        it('should display all signup form elements', () => {
            // check for inputs
            cy.get('.col-7 input').should('have.length', 4)
            cy.contains('Your First Name *').should('be.visible')
            cy.contains('Your Last Name *').should('be.visible')
            cy.contains('Your Email *').should('be.visible')
            cy.contains('Your Password *').should('be.visible')

            // check for button
            cy.contains('sign up').should('be.visible')
        })

        it('should allow typing in all signup inputs', () => {
            cy.get('.col-7 input').eq(0)
                .type('test').should('have.value', 'test')

            cy.get('.col-7 input').eq(1)
                .type('Doe').should('have.value', 'Doe')

            cy.get('.col-7 input').eq(2)
                .type('j@examole.com').should('have.value', 'j@examole.com')

            cy.get('.col-7 input').eq(3)
                .type('password123').should('have.value', 'password123')
        })

        it('should have correct button colors', () => {
            cy.get('.col-5 .q-btn').should('have.class', 'bg-primary')

            cy.get('.col-7 .q-btn').should('have.class', 'bg-positive')
        })
    })

    describe('Form Interactions', () => {
        it('Should handle empty input submissions', () => {
            cy.get('.col-5 button[type="submit"]').click()
            cy.get('.q-notification').should('be.visible').and('contain', 'Email is Required')

            cy.get('.col-7 button[type="submit"]').click()
            cy.get('.q-notification').should('be.visible').and('contain', 'email is Required')
            cy.get('.q-notification').should('be.visible').and('contain', 'password is Required')
            cy.get('.q-notification').should('be.visible').and('contain', 'firstName is Required')
            cy.get('.q-notification').should('be.visible').and('contain', 'lastName is Required')
        })
    })

    describe('Responsive Design', () => {
        it('should maintain layout on different screens', () => {
            cy.viewport(375, 667)
            cy.get('.q-card').should('be.visible')

            cy.viewport(768, 1024)
            cy.get('.col-5').should('be.visible')
            cy.get('.col-7').should('be.visible')
            
            cy.viewport(1200, 800)
            cy.get('.row').should('be.visible')
        })
    })
})

