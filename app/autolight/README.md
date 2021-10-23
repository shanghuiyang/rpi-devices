<img src="../../img/auto-light.gif" width=60% height=60% />

# Auto Light
Auto-Light let you control a led light working with a infared detector together.
the led light will light up when the infared detector detects objects.
and the led will turn off after 30 seconds.

## Connect
infrared detector:
- vcc: 3.3v
- out: GPIO-2
- gnd: GND

led:
- positive: GPIO-16
- negative: GND

```go

          +---------+
          |         |
          | infrared|
          | detector|
          |         |
          +-+--+--+-+
            |  |  |
          gnd out vcc
            |  |  |           +-----------+
            |  |  +-----------+ * 1     o |
            +--|--------------+ * 3     o |
               |              | o       o |
               |              | o       o |         \ | | /
               +--------------+ * 9     o |           ___
                              | o       o |         /     \
                              | o       o |        |-------|
                              | o       o |        |  led  |
                              | o       o |        |       |
                              | o       o |        +--+-+--+
                              | o       o |           | |
                              | o       o |         gnd vcc
                              | o       o |           | |
                              | o       o |           | |
                              | o       o |           | |
                              | o       o |           | |
                              | o    34 * +-----------+ |
                              | o    36 * +-------------+
                              | o       o |
                              | o       o |
                              +-----------+
          
```
