Feature: Test the workshop
  In order to test the workshop 
  I should apply the blueprint 
  and test the resources are created correctly

Scenario: Launch the blueprint
  Given I have a running blueprint
  Then the following resources should be running
    | name                             |
    | resource.network.main            |
    | resource.container.minecraft_web |
    | resource.container.minecraft     |
    | resource.container.vscode        |
    | resource.docs.docs               |