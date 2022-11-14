---
sidebar_position: 5
id: schema_3
title: Schema Resource - Update Method
---

The update method is called by Terraform when it detects there is a difference
between the plan and the state. It gives you an opportunity to mutate a resource.
An example of this could be something like changing metadata for a resource.

Changing metadata does not require a resiource to be recreated it is an
operation that can be updated in place.

For your Minecraft API, this is going to be the easiest update that you need
to make today. Since the schema data can not be modified without recreating
the schema there is no work to do in the `Update` method.

Go ahead and delete all the code inside of the `Update` method, once done
you can progress to the next step, implementing Delete.

