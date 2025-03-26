import EmblaCarousel from 'embla-carousel'

// init carousel
const emblaNode = document.querySelector('.crsl')
const options = { loop: true }
const emblaApi = EmblaCarousel(emblaNode, options)
console.log(emblaApi.slideNodes())