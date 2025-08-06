# CPU NOTES
Se está trabajando con un interprete, eso quiere decir que por cada instrucción, verificaremos que es lo que viene para poder ejecutar las instrucciones como en el Chip-8 y guardarlas en el stack del registro V.

## Fetch-decode loop

Es el main loop para el emulador de la CPU. La CPU real trabaja de esta manera: obtiene un byte o bytes de la memoria los cuales estan almacenados en la posicion apuntada por un registro especial(llamado PC). Luego estos bytes son usados para decidir cual instrucción se debe ejecutra en el CPU. Y cuando decide que debe de hacer, ejecuta la función y lee otro byte o grupo de bytes.

El grupo de bytes que define una sola instrucción en un CPU usualmente son los opcode o operation code.

Para que se ejecute una instrucción se deben pasar diferentes fases usualmente llamadas: fase fetch, fase decode y fase execution.

En la fase fetch se obtiene la data del registro PC y se guarda. Se puede considerar que esta fase como la lectura de código.

En la fase decode se usa la data obtenida anteriormente y se decide que acciones se ejecutarán y se envian señales a las unidades funcionales.

En la fase execution el CPU ejecuta las acciones que se deben realizar para el opcode read.

```
# Ejemplo del libro
while (executed_cycles < cycles_to_execute){
    opcode = memory[PC++];
    instruction = decode(opcode);
    execute(instruction);
}
```


# References

Based on How to write an emulator (CHIP-8 interpreter) 
https://multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/

Opcodes taken from:
https://en.wikipedia.org/wiki/CHIP-8#External_links

and: https://chip8.gulrak.net/

### More documentation:
Emulation Basics: Write your own Chip 8 Emulator/Interpreter: https://omokute.blogspot.com/2012/06/emulation-basics-write-your-own-chip-8.html

Study of the techniques for emulation
programming: http://www.codeslinger.co.uk/files/emu.pdf